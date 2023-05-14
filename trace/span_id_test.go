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

package trace

import (
	"github.com/go-wares/log/adapters"
	"testing"
)

func TestNewSpanId(t *testing.T) {
	s := adapters.NewSpanId()
	t.Logf("span  id: %v", s.Body())
	t.Logf("span str: %v", s.String())
}

func TestNewSpanId2(t *testing.T) {
	// 06ECED2B5CDDDBA1
	// [6 236 237 43 92 221 219 161]
	s := adapters.NewSpanIdFromString("06ECED2B5CDDDBA1")
	t.Logf("span  id: %v", s.Body())
	t.Logf("span str: %v", s.String())
}
