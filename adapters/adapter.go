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
	"github.com/go-wares/log/base"
)

type (
	// LogAdapter
	// 日志适配器.
	LogAdapter interface {
		// Keeper
		// 协程保持.
		Keeper() base.Keeper

		// Send
		// 发送日志.
		Send(line *Line)

		// SetFormatter
		// 设置日志格式.
		SetFormatter(formatter LogFormatter)
	}

	// LogFormatter
	// 日志格式化.
	LogFormatter interface {
		// Byte
		// 转成字符集.
		Byte(line *Line) []byte

		// String
		// 转成字符串.
		String(line *Line) string
	}

	// TraceAdapter
	// 链路适配器.
	TraceAdapter interface {
		// Keeper
		// 协程保持.
		Keeper() base.Keeper

		// Send
		// 发送跨度.
		Send(span Span)
	}
)
