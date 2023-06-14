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

type (
	// TraceId
	// 链路ID接口.
	TraceId interface {
		// Body
		// 获取 byte 列表.
		Body() []byte

		// String
		// 获取字符串.
		//
		// - 长度：16
		String() string
	}

	traceId struct {
		body []byte
		str  string
	}
)

// NewTraceId
// 生成随机链路ID.
func NewTraceId() TraceId {
	o := &traceId{}
	o.body = make([]byte, 16)
	o.str = ID.String(o.body[:])
	return o
}

// NewTraceIdFromString
// 基于字符串反解链路ID.
func NewTraceIdFromString(str string) TraceId {
	o := &traceId{}
	o.str = str
	o.body = ID.Byte(o.str)
	return o
}

// +---------------------------------------------------------------------------+
// | Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *traceId) Body() []byte   { return o.body[:] }
func (o *traceId) String() string { return o.str }

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *traceId) init() *traceId {
	return o
}
