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
	// LogAdapterKafka
	// 消息适配器配置.
	//
	//   # config/log.yaml
	//
	//   log_adapter: kafka
	//   log_adapter_kafka:
	//     address: 127.0.0.1:9092
	//     topic: logs
	LogAdapterKafka struct {
		// 批量阈值.
		// 每次最多批量写入N(默认: 100)条日志.
		Batch int `yaml:"batch" json:"batch"`

		// 保时频率.
		// 每隔固定时长(默认: 350ms)刷盘一次日志.
		Milliseconds int64 `yaml:"milliseconds" json:"milliseconds"`

		// 主机名.
		//
		// - 默认：127.0.0.1:9092
		Host []string

		// 主题名.
		//
		// - 默认：logs
		Topic string
	}
)

func (o *LogAdapterKafka) defaults() {
	if o.Batch == 0 {
		o.Batch = defaultLogAdapterKafkaBatch
	}
	if o.Milliseconds == 0 {
		o.Milliseconds = defaultLogAdapterKafkaMilliseconds
	}
	if len(o.Host) == 0 {
		o.Host = []string{defaultLogAdapterKafkaHost}
	}
	if o.Topic == "" {
		o.Topic = defaultLogAdapterKafkaTopic
	}
}
