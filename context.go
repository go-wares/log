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

package log

import (
	"context"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/config"
)

// Context
// 创建日志上下文.
func Context(parent ...context.Context) context.Context {
	// 1. 根上下文.
	if len(parent) == 0 {
		v := adapters.NewTracing()
		v.TraceId = adapters.NewTraceId()

		return context.WithValue(context.Background(),
			config.OpenTracingKey,
			v,
		)
	}

	// 2. 上级上下文.
	ctx := parent[0]

	// 2.1 子上下文.
	if g := ctx.Value(config.OpenTracingKey); g != nil {
		x := g.(*adapters.Tracing)
		v := adapters.NewTracing()
		v.TraceId = x.TraceId
		v.ParentSpanId = x.SpanId
		return context.WithValue(ctx, config.OpenTracingKey, v)
	}

	// 3. 继承上下文.
	v := adapters.NewTracing()

	// 3.1 隶属主链.
	if g := ctx.Value(config.OpenTracingTraceId); g != nil {
		if str, ok := g.(string); ok && len(str) == 32 {
			v.TraceId = adapters.NewTraceIdFromString(str)
		}
	}

	// 3.2 上级跨度.
	if g := ctx.Value(config.OpenTracingSpanId); g != nil {
		if str, ok := g.(string); ok && len(str) == 16 {
			v.ParentSpanId = adapters.NewSpanIdFromString(str)
		}
	}

	// 3.3 默认链路.
	if v.TraceId == nil {
		v.TraceId = adapters.NewTraceId()
	}

	return context.WithValue(ctx, config.OpenTracingKey, v)
}
