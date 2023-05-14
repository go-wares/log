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
	"encoding/json"
	"github.com/go-wares/log"
	"github.com/go-wares/log/config"
	"github.com/go-wares/log/managers"
	"testing"
	"time"
)

func TestConfig(t *testing.T) {
	t.Logf("[debug:on] %v", config.Config.DebugOn())
	t.Logf("[ info:on] %v", config.Config.InfoOn())
	t.Logf("[ warn:on] %v", config.Config.WarnOn())
	t.Logf("[error:on] %v", config.Config.ErrorOn())
	t.Logf("[fatal:on] %v", config.Config.FatalOn())

	buf, _ := json.MarshalIndent(config.Config, "", "    ")
	t.Logf("config: \n%s", buf)
}

func TestLogAdapter(t *testing.T) {
	c1 := log.Context()

	// managers.Manager.Start()
	defer func() {
		managers.Manager.Stop()
		t.Logf("test log adapter: stopped")
	}()

	time.Sleep(time.Millisecond)

	log.Debugfc(c1, "debug")
	log.Infofc(c1, "info")

	c2 := log.Context(c1)
	log.Warnfc(c2, "warn")
	log.Errorfc(c2, "error")
	log.Fatalfc(c2, "fatal")
}
