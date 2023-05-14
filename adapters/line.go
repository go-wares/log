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
// date: 2023-05-12

package adapters

import (
	"context"
	"fmt"
	"github.com/go-wares/log/base"
	"github.com/go-wares/log/config"
	"sync"
	"time"
)

var (
	linePool sync.Pool
)

type (
	// Line
	// 单行日志.
	Line struct {
		Attr  Attr
		Ctx   context.Context
		Level base.LogLevel
		Text  string
		Time  time.Time

		Tracer               bool
		TraceId              string
		SpanId, ParentSpanId string
	}
)

func NewLine(ctx context.Context, level base.LogLevel, format string, args ...interface{}) *Line {
	if x := linePool.Get(); x != nil {
		return x.(*Line).before(ctx, level, format, args...)
	}

	x := (&Line{}).init()
	return x.before(ctx, level, format, args...)
}

func (o *Line) Release() {
	o.after()
	linePool.Put(o)
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *Line) after() *Line {
	o.Attr = nil
	o.Ctx = nil
	o.Level = base.Off
	o.Text = ""

	if o.Tracer {
		o.Tracer = false
		o.TraceId = ""
		o.SpanId = ""
		o.ParentSpanId = ""
	}

	return o
}

// 初始字段.
func (o *Line) before(ctx context.Context, level base.LogLevel, format string, args ...interface{}) *Line {
	// 1. 必须字段.
	o.Ctx = ctx
	o.Level = level
	o.Time = time.Now()
	o.Text = fmt.Sprintf(format, args...)

	// 2. 堆栈日志.
	if level == base.Fatal {
		o.Text = fmt.Sprintf("%s\n%s", o.Text, Backstack().String())
	}

	// 3. 关联链路.
	if ctx != nil {
		o.openTracing(ctx)
	}

	return o
}

// 构造实例.
func (o *Line) init() *Line { return o }

// 关联链路.
func (o *Line) openTracing(ctx context.Context) {
	// 1. 基于: OpenTelemetry.
	if g := ctx.Value(config.OpenTelemetrySpan); g != nil {
		v := g.(Span)

		o.Tracer = true
		o.TraceId = v.Trace().TraceId().String()
		o.SpanId = v.SpanId().String()

		if p := v.ParentSpanId(); p != nil {
			o.ParentSpanId = p.String()
		}
		return
	}

	// 2. 基于: OpenTracing
	if s := ctx.Value(config.OpenTracingTraceId); s != nil {
		if str, ok := s.(string); ok && len(str) == 32 {
			o.Tracer = true
			o.TraceId = str
		}
	}

	// 3. 跨度参数.
	if o.Tracer {
		// 3.1 跨度ID.
		if s := ctx.Value(config.OpenTracingSpanId); s != nil {
			if str, ok := s.(string); ok && len(str) == 16 {
				o.SpanId = str
			}
		}

		// 3.2 上级跨度.
		if s := ctx.Value(config.OpenTracingParentSpanId); s != nil {
			if str, ok := s.(string); ok && len(str) == 16 {
				o.ParentSpanId = str
			}
		}
	}
}
