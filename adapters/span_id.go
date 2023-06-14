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
	// SpanId
	// 跨度ID接口.
	SpanId interface {
		// Body
		// 获取 byte 列表.
		Body() []byte

		// String
		// 获取字符串.
		//
		// - 长度：16
		String() string
	}

	spanId struct {
		body []byte
		str  string
	}
)

// NewSpanId
// 生成随机跨度.
func NewSpanId() SpanId {
	o := &spanId{}
	o.body = make([]byte, 8)
	o.str = ID.String(o.body[:])
	return o
}

// NewSpanIdFromString
// 基于字符串反解跨度.
func NewSpanIdFromString(str string) SpanId {
	o := &spanId{}
	o.str = str
	o.body = ID.Byte(o.str)
	return o
}

// +---------------------------------------------------------------------------+
// | Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *spanId) Body() []byte   { return o.body[:] }
func (o *spanId) String() string { return o.str }

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *spanId) init() *spanId {
	return o
}
