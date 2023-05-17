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

package log_kafka

import (
	"encoding/json"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/config"
)

type (
	// Data
	// 存储到 Kafka 的数据结构.
	//
	//   {
	//       "time": "2023-05-15 09:10:11.234",
	//       "timestamp_ms": 1684113011234,
	//       "level": "INFO",
	//       "content": "日志内容",
	//
	//       "keywords": {
	//           "id": 1,
	//           "key": "value"
	//       },
	//
	//       "pid": 3721,
	//       "service_addr": "192.168.0.100:8080",
	//       "service_name": "go-wares-log",
	//       "service_version": "1.0"
	//   }
	Data struct {
		Content     string                 `json:"content"`
		Keywords    map[string]interface{} `json:"fields,omitempty"`
		Level       string                 `json:"level"`
		Time        string                 `json:"time"`
		TimestampMs int64                  `json:"timestamp_ms"`

		// +------------------------------------------------------------+
		// | Open tracing                                               |
		// +------------------------------------------------------------+

		ParentSpanId string `json:"parent_span_id,omitempty"`
		SpanId       string `json:"span_id,omitempty"`
		TraceId      string `json:"trace_id,omitempty"`

		// +------------------------------------------------------------+
		// | HTTP Client request                                        |
		// +------------------------------------------------------------+

		RequestMethod string `json:"request_method,omitempty"`
		RequestUrl    string `json:"request_url,omitempty"`
		UserAgent     string `json:"user_agent,omitempty"`

		// +------------------------------------------------------------+
		// | Server fields                                              |
		// +------------------------------------------------------------+

		// 进程号
		Pid int `json:"pid"`

		// 服务地址.
		//
		// 例如：192.168.0.100
		//      192.168.0.100:8080
		ServiceAddr    []string `json:"service_addr,omitempty"`
		ServiceName    string   `json:"service_name"`
		ServiceVersion string   `json:"service_version"`
	}

	// Formatter
	// 格式化.
	Formatter struct{}
)

// Byte
// 转成Byte字符集.
func (o *Formatter) Byte(line *adapters.Line) (body []byte) {
	v := &Data{
		Content:        line.Text,
		Level:          line.Level.String(),
		TimestampMs:    line.Time.UnixMilli(),
		Time:           line.Time.Format("2006-01-02T15:04:05.999Z"),
		Pid:            config.Config.Pid,
		ServiceAddr:    config.Config.Addr,
		ServiceName:    config.Config.Name,
		ServiceVersion: config.Config.Version,
	}

	// 关键字段.
	if line.Attr.Count() > 0 {
		v.Keywords = line.Attr
	}

	// 调用链路.
	if line.Tracer {
		v.ParentSpanId = line.ParentSpanId
		v.SpanId = line.SpanId
		v.TraceId = line.TraceId
	}

	if buf, err := json.Marshal(v); err == nil {
		return buf
	}

	return nil
}

// String
// 转成字符串.
func (o *Formatter) String(line *adapters.Line) string {
	if buf := o.Byte(line); buf != nil {
		return string(buf)
	}
	return ""
}

func (o *Formatter) init() *Formatter { return o }
