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

package log_file

import (
	"fmt"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/config"
)

type (
	Formatter struct{}
)

// String
// 转成字符串.
func (o *Formatter) String(line *adapters.Line) string {
	// 1. 基础日志.
	text := fmt.Sprintf("[%s][%s]", line.Time.Format(config.Config.LogTimeFormat), line.Level)

	// 2. 日志正文.
	if line.Attr != nil {
		text = fmt.Sprintf("%s %s - %s", text, line.Attr.Json(), line.Text)
	} else {
		text = fmt.Sprintf("%s %s", text, line.Text)
	}

	// 3. 单行日志.
	return text
}

func (o *Formatter) init() *Formatter {
	return o
}
