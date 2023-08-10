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
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/config"
	"net/http"
)

type (
	// 链路结构.
	trace struct {
		ctx          context.Context
		name         string
		parentSpanId adapters.SpanId
		traceId      adapters.TraceId
	}
)

// NewTrace
// 创建根链路.
func NewTrace(name string) adapters.Trace {
	o := (&trace{name: name}).init()
	o.traceId = adapters.NewTraceId()
	o.ctx = context.WithValue(context.Background(), config.OpenTelemetryTrace, o)
	return o
}

// NewTraceFromContext
// 基于上下文创建链路.
func NewTraceFromContext(ctx context.Context, name string) adapters.Trace {
	// 1. 链路复用.
	if g := ctx.Value(config.OpenTelemetryTrace); g != nil {
		if o, ok := g.(*trace); ok {
			return o
		}
	}

	// 2. 创建链路
	o := (&trace{name: name}).init()

	// 3. 复用链路ID.
	if g := ctx.Value(config.OpenTracingTraceId); g != nil {
		if str, ok := g.(string); ok && len(str) == 32 {
			o.traceId = adapters.NewTraceIdFromString(str)
		}
	}

	// 4. 随机链路ID.
	if o.traceId == nil {
		o.traceId = adapters.NewTraceId()
	}

	// 5. 上级跨度.
	if g := ctx.Value(config.OpenTracingSpanId); g != nil {
		if str, ok := g.(string); ok && len(str) == 16 {
			o.parentSpanId = adapters.NewSpanIdFromString(str)
		}
	}

	// 6. 设置上下文.
	o.ctx = context.WithValue(context.Background(), config.OpenTelemetryTrace, o)
	return o
}

func NewTraceFromRequest(req *http.Request, name string) adapters.Trace {
	o := (&trace{name: name}).init()

	if s := req.Header.Get(config.OpenTracingTraceId); s != "" {
		o.traceId = adapters.NewTraceIdFromString(s)
	} else {
		o.traceId = adapters.NewTraceId()
		req.Header.Set(config.OpenTracingTraceId, o.traceId.String())
	}

	if s := req.Header.Get(config.OpenTracingSpanId); s != "" {
		o.parentSpanId = adapters.NewSpanIdFromString(s)
	} else {
		o.parentSpanId = adapters.NewSpanId()
		req.Header.Set(config.OpenTracingSpanId, o.parentSpanId.String())
	}

	return o
}

func (o *trace) Begin(name string) adapters.Span {
	v := spanPool.Get().(*span).before()
	v.name = name
	v.parentSpanId = o.parentSpanId
	v.trace = o
	v.withCtx(o.ctx)
	return v
}

// +---------------------------------------------------------------------------+
// | Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *trace) Context() context.Context  { return o.ctx }
func (o *trace) Name() string              { return o.name }
func (o *trace) TraceId() adapters.TraceId { return o.traceId }

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *trace) init() *trace {
	return o
}
