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

package tests

import (
	"github.com/go-wares/log"
	"github.com/go-wares/log/managers"
	"testing"
	"time"
)

func TestField(t *testing.T) {
	defer func() {
		managers.Manager.Stop()
		t.Logf("test log adapter: stopped")
	}()

	time.Sleep(time.Millisecond)

	field := log.Field{"key": "value", "uid": 1}
	field.Debug("debug")
	field.Info("info")
	field.Warn("warn")
	field.Error("error")
	field.Fatal("fatal")
}
