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

package managers

import (
	"context"
	"fmt"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/adapters/log_file"
	"github.com/go-wares/log/adapters/log_term"
	"github.com/go-wares/log/adapters/trace_jaeger"
	"github.com/go-wares/log/base"
	"github.com/go-wares/log/config"
	"github.com/go-wares/log/trace"
	"net"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
)

var (
	// Manager
	// 管理器实例.
	Manager Management
)

type (
	// Management
	// 基础管理器.
	Management interface {
		GetLogAdapter() adapters.LogAdapter
		Log(ctx context.Context, fields map[string]interface{}, level base.LogLevel, format string, args ...interface{})
		Start()
		Stop()
	}

	manager struct {
		cancel context.CancelFunc
		ctx    context.Context
		keeper base.Keeper
		mu     sync.RWMutex
		name   string

		logAdapter   adapters.LogAdapter
		traceAdapter adapters.TraceAdapter
	}
)

func (o *manager) GetLogAdapter() adapters.LogAdapter { return o.logAdapter }

func (o *manager) Log(ctx context.Context, fields map[string]interface{}, level base.LogLevel, format string, args ...interface{}) {
	if o.logAdapter != nil {
		line := adapters.NewLine(ctx, level, format, args...)

		if fields != nil {
			line.Attr = fields
		}

		o.logAdapter.Send(line)
	}
}

func (o *manager) Start() {
	o.mu.Lock()

	if o.ctx != nil {
		o.mu.Unlock()
		return
	}

	o.ctx, o.cancel = context.WithCancel(context.Background())
	o.mu.Unlock()

	defer func() {
		o.mu.Lock()
		defer o.mu.Unlock()

		if o.ctx.Err() == nil {
			o.cancel()
		}

		o.cancel = nil
		o.ctx = nil
	}()

	if err := o.keeper.Start(o.ctx); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v", err)
	}
}

func (o *manager) Stop() {
	o.mu.RLock()
	if o.ctx != nil && o.ctx.Err() == nil {
		o.cancel()
	}
	o.mu.RUnlock()

	for {
		if o.keeper.Stopped() {
			return
		}

		time.Sleep(time.Millisecond * 10)
	}
}

// +---------------------------------------------------------------------------+
// | Event methods                                                             |
// +---------------------------------------------------------------------------+

func (o *manager) onListen(ctx context.Context) (ignored bool) {
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

func (o *manager) init() *manager {
	o.name = fmt.Sprintf("log-manager")
	o.keeper = base.NewKeeper(o.name).Listen(o.onListen)

	o.initLogAdapter()
	o.initTraceAdapter()
	o.initAdapterResource()
	return o
}

func (o *manager) initAdapterResource() {
	// 架构名称.
	adapters.Resource.
		Set("deploy.arch", fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)).
		Set("deploy.go", runtime.Version()).
		Set("deploy.pid", os.Getpid())

	// 主机名称.
	if s, se := os.Hostname(); se == nil {
		adapters.Resource.Set("deploy.host", s)
	}

	// 主机地址.
	if l, le := net.InterfaceAddrs(); le == nil {
		ls := make([]string, 0)
		for _, la := range l {
			if ipn, ok := la.(*net.IPNet); ok && !ipn.IP.IsLoopback() {
				if ipn.IP.To4() != nil {
					ls = append(ls, ipn.IP.String())
				}
			}
		}
		adapters.Resource.Set("deploy.addr", strings.Join(ls, ", "))
	}
}

func (o *manager) initLogAdapter() {
	// 1. 日志适配器.
	switch config.Config.LogAdapter {
	case base.File:
		o.logAdapter = log_file.New()
	case base.Term:
		o.logAdapter = log_term.New()
	}

	// 2. 加为子 Keeper.
	if o.logAdapter != nil {
		o.keeper.Add(o.logAdapter.Keeper())
	}
}

func (o *manager) initTraceAdapter() {
	// 1. 链路适配器.
	switch config.Config.TraceAdapter {
	case base.Jaeger:
		o.traceAdapter = trace_jaeger.New()
	}

	// 2. 加为子 Keeper.
	if o.traceAdapter != nil {
		trace.SpanPublish = o.traceAdapter.Send
		o.keeper.Add(o.traceAdapter.Keeper())
	}
}
