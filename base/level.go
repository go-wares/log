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

package base

import (
	"strings"
)

type (
	// Level
	// 级别名称.
	Level string

	// LogLevel
	// 级别编码.
	LogLevel int
)

func (x Level) LogLevel() (l Level, ll LogLevel) {
	var (
		s  = strings.ToUpper(string(x))
		ls string
	)

	// 未配置.
	if s == "" {
		l = Level(LogLevelText[Info])
		ll = Info
		return
	}

	// 有效配置.
	for ll, ls = range LogLevelText {
		if ls == s {
			l = Level(s)
			return
		}
	}

	ll = Off
	return
}

const (
	Off LogLevel = iota
	Fatal
	Error
	Warn
	Info
	Debug
)

var LogLevelText = map[LogLevel]string{
	Fatal: "FATAL",
	Error: "ERROR",
	Warn:  "WARN",
	Info:  "INFO",
	Debug: "DEBUG",
}

func (n LogLevel) String() string {
	return LogLevelText[n]
}
