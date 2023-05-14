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

package adapters

import (
	"encoding/json"
)

var (
	// Resource
	// 系统资源.
	//
	//   {
	//       "host": "127.0.0.1"
	//   }
	Resource = make(Attr)
)

type (
	Attr map[string]interface{}
)

func (o Attr) Set(key string, value interface{}) Attr {
	if o != nil {
		o[key] = value
	}
	return o
}

func (o Attr) Json() string {
	buf, _ := json.Marshal(o)
	return string(buf)
}
