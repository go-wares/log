// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// author: wsfuyibing <websearch@163.com>
// date: 2023-05-14

package trace

import (
	"context"
	"encoding/json"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/base"
	"github.com/go-wares/log/config"
	"net/http"
	"sync"
	"time"
)

var (
	spanPool    sync.Pool
	spanNilTime = time.Unix(0, 0)
)

type (
	// 跨度结构.
	span struct {
		attr                 adapters.Attr
		ctx                  context.Context
		endTime, startTime   time.Time
		lines                []*adapters.Line
		mu                   *sync.RWMutex
		name                 string
		running              bool
		spanId, parentSpanId adapters.SpanId
		trace                adapters.Trace
	}
)

// NewSpan
// 创建跨度.
func NewSpan(name string) adapters.Span {
	return NewTrace(name).Begin(name)
}

// NewSpanFromContext
// 基于上下文创建跨度.
func NewSpanFromContext(ctx context.Context, name string) adapters.Span {
	// 1. 基于跨度.
	//    从跨度(Span)上开启子跨度(Span).
	if x := ctx.Value(config.OpenTelemetrySpan); x != nil {
		if o, ok := x.(adapters.Span); ok {
			return o.Child(name)
		}
	}

	// 2. 创建跨度.
	return NewTraceFromContext(ctx, name).Begin(name)
}

// NewSpanFromRequest
// 基于HTTP请求创建跨度.
func NewSpanFromRequest(req *http.Request, name string) adapters.Span {
	o := NewTraceFromRequest(req, name)
	s := o.Begin(name)
	s.(*span).ReadRequest(req)
	return s
	// o := NewSpanFromContext(req.Context(), name)
	// o.(*span).ReadRequest(req)
	// return o
}

func SpanExists(ctx context.Context) (span adapters.Span, exists bool) {
	if ctx != nil {
		if v := ctx.Value(config.OpenTelemetrySpan); v != nil {
			return v.(adapters.Span), true
		}
	}
	return nil, false
}

func (o *span) Child(name string) adapters.Span {
	v := spanPool.Get().(*span).before()
	v.name = name
	v.parentSpanId = o.spanId
	v.trace = o.trace
	v.withCtx(o.ctx)
	return v
}

// +---------------------------------------------------------------------------+
// | Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *span) Attr() adapters.Attr           { return o.attr }
func (o *span) Context() context.Context      { return o.ctx }
func (o *span) End()                          { o.end() }
func (o *span) EndTime() time.Time            { return o.endTime }
func (o *span) Logs() []*adapters.Line        { return o.lines }
func (o *span) Name() string                  { return o.name }
func (o *span) ParentSpanId() adapters.SpanId { return o.parentSpanId }
func (o *span) Release()                      { o.after(); spanPool.Put(o) }
func (o *span) SpanId() adapters.SpanId       { return o.spanId }
func (o *span) StartTime() time.Time          { return o.startTime }
func (o *span) Trace() adapters.Trace         { return o.trace }

// +---------------------------------------------------------------------------+
// | Span logs                                                                 |
// +---------------------------------------------------------------------------+

func (o *span) Debug(format string, args ...interface{}) {
	if config.Config.DebugOn() {
		o.log(base.Debug, format, args...)
	}
}

func (o *span) Info(format string, args ...interface{}) {
	if config.Config.InfoOn() {
		o.log(base.Info, format, args...)
	}
}

func (o *span) Warn(format string, args ...interface{}) {
	if config.Config.WarnOn() {
		o.log(base.Warn, format, args...)
	}
}

func (o *span) Error(format string, args ...interface{}) {
	if config.Config.ErrorOn() {
		o.log(base.Error, format, args...)
	}
}

func (o *span) Fatal(format string, args ...interface{}) {
	if config.Config.FatalOn() {
		o.log(base.Fatal, format, args...)
	}
}

// +---------------------------------------------------------------------------+
// | Request methods                                                           |
// +---------------------------------------------------------------------------+

func (o *span) ReadRequest(request *http.Request) {
	buf, _ := json.Marshal(request.Header)

	o.attr.Set("http.header", string(buf))
	o.attr.Set("http.protocol", request.Proto)
	o.attr.Set("http.request.method", request.Method)
	o.attr.Set("http.request.uri", request.URL.Path)
	o.attr.Set("http.user.agent", request.UserAgent())
}

func (o *span) WriteRequest(request *http.Request) {
	request.Header.Set(config.OpenTracingTraceId, o.trace.TraceId().String())
	request.Header.Set(config.OpenTracingSpanId, o.spanId.String())
	request.Header.Set(config.OpenTracingSampled, config.OpenTracingSampledFlag)
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

// 后置跨度.
func (o *span) after() {
	// 1. 释放日志.
	for _, line := range o.lines {
		line.Release()
	}

	// 2. 重置字段.
	o.attr = nil
	o.ctx = nil
	o.endTime = spanNilTime
	o.lines = nil
	o.parentSpanId = nil
	o.spanId = nil
	o.trace = nil
}

// 前置跨度.
func (o *span) before() *span {
	o.attr = make(adapters.Attr)
	o.lines = make([]*adapters.Line, 0)
	o.running = true
	o.spanId = adapters.NewSpanId()
	o.startTime = time.Now()
	return o
}

// 结束跨度.
func (o *span) end() {
	o.mu.Lock()

	// 1. 已经结束.
	if !o.running {
		o.mu.Unlock()
		return
	}

	// 2. 结束跨度.
	o.endTime = time.Now()
	o.running = false
	o.mu.Unlock()

	// 3. 上报链路.
	if TraceManager != nil {
		TraceManager.Send(o)
	} else {
		o.Release()
	}
}

// 构造跨度.
func (o *span) init() *span {
	o.mu = &sync.RWMutex{}
	return o
}

// 记录日志.
func (o *span) log(level base.LogLevel, format string, args ...interface{}) {
	// 1. 跨度日志.
	func() {
		// 加入列表.
		o.mu.Lock()
		defer o.mu.Unlock()

		o.lines = append(o.lines, adapters.NewLine(
			nil, level, format, args...,
		))
	}()

	// 2. 日志同步.
	//    当记录链路(跨度)日志时, 同步写一份到日志系统中.
	if *config.Config.TraceAdapterSyncLog {
		line := adapters.NewLine(o.ctx, level, format, args...)
		// line.Attr = o.Attr()
		LogManager.Send(line)
	}
}

func (o *span) withCtx(ctx context.Context) {
	o.ctx = context.WithValue(ctx, config.OpenTelemetrySpan, o)
}
