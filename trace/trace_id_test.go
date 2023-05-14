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

func TestNewTraceId(t *testing.T) {
	s := adapters.NewTraceId()
	t.Logf("span  id: %v", s.Body())
	t.Logf("span str: %v", s.String())
}

func TestNewTraceId2(t *testing.T) {
	// 0c80b3e216a8341155daf246198bdf37
	// [12 128 179 226 22 168 52 17 85 218 242 70 25 139 223 55]
	s := adapters.NewTraceIdFromString("0c80b3e216a8341155daf246198bdf37")
	t.Logf("span  id: %v", s.Body())
	t.Logf("span str: %v", s.String())
}
