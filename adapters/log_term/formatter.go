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

package log_term

import (
	"fmt"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/base"
	"github.com/go-wares/log/config"
)

type (
	// Formatter
	// 格式化.
	Formatter struct{}
)

// Byte
// 转成Byte字符集.
func (o *Formatter) Byte(_ *adapters.Line) []byte { return nil }

// String
// 转成字符串.
func (o *Formatter) String(line *adapters.Line) string {
	var (
		// 日志正文.
		text = fmt.Sprintf("[%s][%s]",
			line.Time.Format(config.Config.LogTimeFormat),
			line.Level,
		)
	)

	// 1. 链路信息.
	if line.Tracer {
		text = fmt.Sprintf("%s [span-id=%s]",
			text,
			line.SpanId,
		)
	}

	// 2. 绑定字段.
	if len(line.Attr) > 0 {
		text = fmt.Sprintf("%s %s",
			text,
			line.Attr.Json(),
		)
	}

	// 3. 用户正文.
	text = fmt.Sprintf("%s %s",
		text,
		line.Text,
	)

	// 4. 日志着色.
	if *config.Config.LogAdapterTerm.Color {
		return o.color(line.Level, text)
	}

	return text
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

// 着色.
func (o *Formatter) color(level base.LogLevel, str string) string {
	// 黄色红背.
	if level == base.Fatal {
		return fmt.Sprintf("%c[%d;%d;%dm%s%c[0m",
			0x1B, 0, 43, 31, str, 0x1B,
		)
	}

	// 红色 - 1
	if level == base.Error {
		return fmt.Sprintf("%c[%dm%s%c[0m",
			0x1B, 31, str, 0x1B,
		)
	}

	// 黄色 - 3.
	if level == base.Warn {
		return fmt.Sprintf("%c[%dm%s%c[0m",
			0x1B, 33, str, 0x1B,
		)
	}

	// 蓝色 - 4
	if level == base.Info {
		return fmt.Sprintf("%c[%dm%s%c[0m",
			0x1B, 34, str, 0x1B,
		)
	}

	// 灰色 - 7.
	if level == base.Debug {
		return fmt.Sprintf("%c[%dm%s%c[0m",
			0x1B, 37, str, 0x1B,
		)
	}

	return str
}

func (o *Formatter) init() *Formatter {
	return o
}
