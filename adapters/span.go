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

package adapters

import (
	"context"
	"time"
)

type (
	Span interface {
		// Attr
		// 跨度属性.
		Attr() Attr

		// Child
		// 新建子跨度.
		Child(name string) Span

		// Context
		// 跨度上下文.
		Context() context.Context

		// End
		// 跨度结束.
		//
		// 当此方法被显现调用后, 上报链接跨度到指定服务上.
		End()

		// EndTime
		// 结束时间.
		EndTime() time.Time

		Logs() []*Line

		Name() string

		ParentSpanId() SpanId

		// Release
		// 释放回池.
		Release()

		SpanId() SpanId

		// StartTime
		// 开始时间.
		StartTime() time.Time

		Trace() Trace

		Debug(format string, args ...interface{})
		Info(format string, args ...interface{})
		Warn(format string, args ...interface{})
		Error(format string, args ...interface{})
		Fatal(format string, args ...interface{})
	}
)
