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
// date: 2023-05-13

package config

type (
	// LogAdapterTerm
	// 终端适配器配置.
	//
	//   # config/log.yaml
	//
	//   log_adapter: term
	//   log_adapter_term:
	//     color: false
	LogAdapterTerm struct {
		// 着色.
		// 打印到终端上的日志是否包含颜色.
		Color *bool
	}
)

func (o *LogAdapterTerm) defaults(_ *Configuration) {
	if o.Color == nil {
		o.Color = &defaultLogAdapterTermColor
	}
}
