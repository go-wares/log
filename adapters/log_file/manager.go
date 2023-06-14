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

package log_file

import (
	"context"
	"fmt"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/base"
	"github.com/go-wares/log/config"
	"os"
	"sync"
	"time"
)

type (
	// Manager
	// 日志管理器.
	//
	// 发送用户日志到文件中, 记录在运行节点的磁盘上.
	Manager struct {
		bucket      *adapters.Bucket
		directories map[string]bool
		formatter   adapters.LogFormatter
		keeper      base.Keeper
		mu          sync.RWMutex
		name        string
	}
)

func New() adapters.LogAdapter {
	return (&Manager{}).init()
}

func (o *Manager) Keeper() base.Keeper { return o.keeper }

// Send
// 加入数据桶.
//
// 若数据桶积压数量超过指定值时, 立即刷盘保存.
func (o *Manager) Send(line *adapters.Line) {
	if n := o.bucket.Add(line); n >= config.Config.LogAdapterFile.Batch {
		go o.save()
	}
}

// SetFormatter
// 设置格式.
func (o *Manager) SetFormatter(formatter adapters.LogFormatter) {
	o.formatter = formatter
}

// +---------------------------------------------------------------------------+
// | Event methods                                                             |
// +---------------------------------------------------------------------------+

func (o *Manager) onAfter(ctx context.Context) (ignored bool) {
	if o.bucket.Count() > 0 {
		o.save()
		return o.onAfter(ctx)
	}
	return
}

func (o *Manager) onListen(ctx context.Context) (ignored bool) {
	// 1. 定时保存.
	//    每隔指定时长(默认: 350ms)上报一次日志.
	ticker := time.NewTicker(time.Duration(config.Config.LogAdapterFile.Milliseconds) * time.Millisecond)

	// 2. 关闭定时.
	defer ticker.Stop()

	// 3. 监听信号.
	for {
		select {
		case <-ticker.C:
			go o.save()
		case <-ctx.Done():
			return
		}
	}
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *Manager) init() *Manager {
	o.bucket = adapters.NewBucket()
	o.directories = make(map[string]bool)
	o.formatter = (&Formatter{}).init()
	o.name = fmt.Sprintf("log-file-manager")
	o.keeper = base.NewKeeper(o.name).
		After(o.onAfter).
		Listen(o.onListen)
	return o
}

func (o *Manager) mkdir(path string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	// 已经创建.
	if _, ok := o.directories[path]; ok {
		return
	}

	// 创建目录.
	o.directories[path] = true
	if err := os.MkdirAll(path, os.ModePerm); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "make dir: %v\n", err)
	}
}

func (o *Manager) save() {
	var (
		list, count = o.bucket.Popn(config.Config.LogAdapterFile.Batch)
		writer      *Writer
	)

	// 1. 空数据桶.
	if count == 0 {
		return
	}

	// 2. 释放实例.
	defer func() {
		// 2.1 释放日志.
		for _, v := range list {
			v.(*adapters.Line).Release()
		}

		// 2.2 释放实例.
		if writer != nil {
			writer.Release()
		}
	}()

	// 3. 获取实例.
	writer = NewWriter()
	writer.Send(o, list)
}
