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
// date: 2023-04-18

package log

import (
	"github.com/go-wares/log/config"
	"github.com/go-wares/log/managers"
	"sync"
	"time"
)

func init() {
	new(sync.Once).Do(func() {
		if *config.Config.AutoStart {
			go managers.Manager.Start()
			time.Sleep(time.Millisecond)
		}
	})
}
