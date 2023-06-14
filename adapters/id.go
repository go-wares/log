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

import (
	cr "crypto/rand"
	eb "encoding/binary"
	"encoding/hex"
	mr "math/rand"
	"sync"
)

var (
	// ID
	// 生成ID实例.
	ID Identify
)

type (
	// Identify
	// 生成链路/跨度ID接口.
	Identify interface {
		// Byte
		// 从 String 转成 Byte.
		Byte(str string) (body []byte)

		// String
		// 从 Body 转成 String.
		String(body []byte) (str string)
	}

	id struct {
		sync.Mutex
		data   int64
		err    error
		random *mr.Rand
	}
)

// +---------------------------------------------------------------------------+
// | Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *id) Byte(str string) []byte {
	body := make([]byte, len(str)/2)
	if d, de := hex.DecodeString(str); de == nil {
		copy(body[:], d)
	}
	return body
}

func (o *id) String(s []byte) string {
	o.Lock()
	defer o.Unlock()

	o.random.Read(s[:])
	return hex.EncodeToString(s[:])
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *id) init() *id {
	o.err = eb.Read(cr.Reader, eb.LittleEndian, &o.data)
	o.random = mr.New(mr.NewSource(o.data))
	return o
}
