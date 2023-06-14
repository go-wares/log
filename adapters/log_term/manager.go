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

package log_term

import (
	"context"
	"fmt"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/base"
	"os"
)

type (
	// Manager
	// 终端管理器.
	//
	// 发送用户日志到终端(Terminal)上打印.
	Manager struct {
		formatter adapters.LogFormatter
		keeper    base.Keeper
		name      string
	}
)

func New() adapters.LogAdapter {
	return (&Manager{}).
		init()
}

func (o *Manager) Keeper() base.Keeper { return o.keeper }

func (o *Manager) Send(line *adapters.Line) {
	defer line.Release()
	_, _ = fmt.Fprintf(os.Stdout, fmt.Sprintf("%s\n", o.formatter.String(line)))
}

func (o *Manager) SetFormatter(formatter adapters.LogFormatter) {
	o.formatter = formatter
}

// +---------------------------------------------------------------------------+
// | Event methods                                                             |
// +---------------------------------------------------------------------------+

func (o *Manager) onListen(ctx context.Context) (ignored bool) {
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *Manager) init() *Manager {
	o.formatter = (&Formatter{}).init()
	o.name = fmt.Sprintf("log-term-manager")
	o.keeper = base.NewKeeper(o.name).Listen(o.onListen)
	return o
}
