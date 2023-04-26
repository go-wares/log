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
// date: 2023-04-18

package log

import (
	"fmt"
	"github.com/go-wares/log/config"
	"os"
)

func Debug(format string, args ...interface{}) {
	term(config.Debug, format, args)
}

func Error(format string, args ...interface{}) {
	term(config.Error, format, args)
}

func Fatal(format string, args ...interface{}) {
	term(config.Fatal, format, args)
}

func Info(format string, args ...interface{}) {
	term(config.Info, format, args)
}

func Warn(format string, args ...interface{}) {
	term(config.Warn, format, args)
}

func term(level config.Level, format string, args []interface{}) {
	str := fmt.Sprintf("[%s] %s", level, fmt.Sprintf(format, args...))

	switch level {
	case config.Fatal:
		str = fmt.Sprintf("%c[%d;%d;%dm%s%c[0m",
			0x1B, 0, 41, 33, str, 0x1B,
		)
	case config.Error:
		str = fmt.Sprintf("%c[%d;%d;%dm%s%c[0m",
			0x1B, 0, 0, 31, str, 0x1B,
		)
	case config.Warn:
		str = fmt.Sprintf("%c[%d;%d;%dm%s%c[0m",
			0x1B, 0, 0, 33, str, 0x1B,
		)
	case config.Info:
		str = fmt.Sprintf("%c[%d;%d;%dm%s%c[0m",
			0x1B, 0, 0, 34, str, 0x1B,
		)
	default:
	}

	_, _ = fmt.Fprintln(os.Stdout, str)
}
