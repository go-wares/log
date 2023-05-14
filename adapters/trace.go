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
)

type (
	Trace interface {
		// Begin
		// 开启跨度.
		Begin(name string) Span

		// Context
		// 获取链路上下文.
		Context() context.Context

		// Name
		// 链路名称.
		Name() string

		// TraceId
		// 获取链路ID.
		TraceId() TraceId
	}
)
