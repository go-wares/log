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

package trace_jaeger

import (
	"context"
	"fmt"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/base"
	"github.com/go-wares/log/config"
	"time"
)

type (
	// Manager
	// 链路(Jaeger)管理器.
	Manager struct {
		bucket    *adapters.Bucket
		formatter *formatter
		keeper    base.Keeper
		name      string
	}
)

func New() adapters.TraceAdapter {
	return (&Manager{}).init()
}

// +---------------------------------------------------------------------------+
// | Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *Manager) Keeper() base.Keeper { return o.keeper }

func (o *Manager) Send(span adapters.Span) {
	if n := o.bucket.Add(span); n >= config.Config.TraceAdapterJaeger.Batch {
		go o.save()
	}
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
	//    每隔指定时长(默认: 350ms)上报一次链路跨度.
	ticker := time.NewTicker(time.Duration(config.Config.TraceAdapterJaeger.Milliseconds) * time.Millisecond)

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
	o.formatter = (&formatter{}).init()
	o.name = fmt.Sprintf("trace-jaeger-manager")
	o.keeper = base.NewKeeper(o.name).
		After(o.onAfter).
		Listen(o.onListen)
	return o
}

func (o *Manager) save() {
	var (
		buf, count = o.bucket.Popn(config.Config.TraceAdapterJaeger.Batch)
		list       = make([]adapters.Span, 0)
	)

	// 1. 空数据桶.
	if count == 0 {
		return
	}

	// 2. 释放实例.
	defer func() {
		// 2.1 释放日志.
		for _, v := range list {
			v.Release()
		}
	}()

	// 3. 获取实例.
	for _, x := range buf {
		list = append(list, x.(adapters.Span))
	}

	v := NewWriter()
	v.Send(o.formatter, list...)
	v.Release()
}
