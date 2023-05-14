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

package config

type (
	// TraceAdapterJaeger
	// Jaeger 链路配置.
	TraceAdapterJaeger struct {
		Endpoint string `yaml:"endpoint" json:"endpoint"`
		Username string `yaml:"username" json:"username"`
		Password string `yaml:"password" json:"password"`

		Batch        int    `yaml:"batch" json:"batch"`
		Milliseconds int    `yaml:"milliseconds" json:"milliseconds"`
		Topic        string `yaml:"topic" json:"topic"`
	}
)

func (o *TraceAdapterJaeger) defaults() {
	if o.Batch == 0 {
		o.Batch = defaultTraceAdapterJaegerBatch
	}
	if o.Milliseconds == 0 {
		o.Milliseconds = defaultTraceAdapterJaegerMilliseconds
	}
}
