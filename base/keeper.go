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

package base

import (
	"context"
	"fmt"
	"os"
	"sync"
)

type (
	// Keeper
	// 协程保持.
	//
	// 工作过程模拟进程, 每个 Keeper 独占一个协程, 当 Keeper 退出或重启时, 先等待
	// 子 Keeper 成功退出.
	Keeper interface {
		// Add
		// 添加子 Keeper.
		Add(keeper Keeper) (yes bool)

		// Del
		// 删除子 Keeper.
		Del(v Keeper) (yes bool)

		// DelAuto
		// 自动删除.
		//
		// 当值为 true 时, 若线程退出则通知 parent Keeper 删除子 Keeper.
		DelAuto(ad bool) Keeper

		// After
		// 注册后置执行器.
		//
		// 仅执行1次, 当过程执行器执行完成且子Keeper已经退出后触发.
		After(handlers ...KeeperHandler) Keeper

		// Before
		// 注册前置执行器.
		//
		// 仅执行1次, 若任一 KeeperHandler 返回 true (ignored) 时, 跳过 After 和 Listen
		// 注册过的执行器, 并退出协程.
		Before(handlers ...KeeperHandler) Keeper

		// Listen
		// 注册过程执行器.
		//
		// 可执行1+次, 当 Keeper 启动或重启时, 都会触发1次.
		Listen(handlers ...KeeperHandler) Keeper

		// Name
		// 返回名称.
		Name() string

		// Panic
		// 注册异常执行器.
		//
		// 当在 KeeperHandler 出现运行异常时自动触发.
		Panic(handler KeeperPanicHandler) Keeper

		// Restart
		// 重启 Keeper.
		Restart()

		SetParent(p Keeper)

		// Start
		// 启动 Keeper 协程.
		Start(ctx context.Context) (err error)

		// Stop
		// 退出 Keep 协程.
		Stop()

		// Stopped
		// 已退出状态.
		Stopped() bool
	}

	KeeperHandler      func(ctx context.Context) (ignored bool)
	KeeperPanicHandler func(ctx context.Context, v interface{})

	// 协程保持.
	keeper struct {
		ctx    context.Context
		cancel context.CancelFunc

		children                          map[string]Keeper
		listAfter, listBefore, listListen []KeeperHandler
		panicHandler                      KeeperPanicHandler

		mu                           sync.RWMutex
		name                         string
		parent                       Keeper
		started, restart, autoDelete bool
	}
)

func NewKeeper(name string) Keeper {
	return (&keeper{
		name: name,
	}).init()
}

// +---------------------------------------------------------------------------+
// | Interface methods                                                         |
// +---------------------------------------------------------------------------+

func (o *keeper) Name() string {
	return o.name
}

func (o *keeper) Restart() {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.ctx != nil && o.ctx.Err() == nil {
		o.restart = true
		o.cancel()
	}
}

func (o *keeper) Start(ctx context.Context) (err error) {
	o.mu.Lock()

	// 1. 重复启动.
	if o.started {
		o.mu.Unlock()
		err = fmt.Errorf("%s started already", o.name)
		return
	}

	// 2. 取消启动.
	//    若上下文为 nil 或 上下文已经退出时, 取消启动.
	if ctx == nil || ctx.Err() != nil {
		o.mu.Unlock()
		err = fmt.Errorf("%s context is cancelled", o.name)
		return
	}

	// 3. 锁定状态.
	o.started = true
	o.mu.Unlock()

	// 4. 退出协程.
	defer func() {
		if o.autoDelete && o.parent != nil {
			o.parent.Del(o)
		}

		o.reset()
	}()

	// 5. 前置执行器.
	if o.runHandlers(ctx, o.listBefore) {
		return
	}

	// 6. 后置执行器.
	defer o.runHandlers(ctx, o.listAfter)

	// 7. 过程执行器.
	for {
		// 7.1 状态检测.
		if func() bool {
			o.mu.Lock()
			defer o.mu.Unlock()
			if o.restart {
				o.restart = false
				return false
			}
			return true
		}() {
			return
		}

		// 7.2 上下文退出.
		if ctx.Err() != nil {
			return
		}

		// 7.3 过程上下文.
		o.ctx, o.cancel = context.WithCancel(ctx)

		// 7.3.1 启动子 Keeper.
		o.childrenStart(o.ctx)

		// 7.3.2 启动过程执行器.
		o.runHandlers(o.ctx, o.listListen)

		// 7.4 过程执行器结束.
		if o.ctx.Err() == nil {
			o.cancel()
		}

		// 7.4.1 等子 Keeper 退出.
		o.childrenWaiter()
	}
}

func (o *keeper) Stop() {
	o.mu.Lock()
	defer o.mu.Unlock()

	if o.ctx != nil && o.ctx.Err() == nil {
		o.restart = true
		o.cancel()
	}
}

func (o *keeper) Stopped() bool {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return !o.started
}

// +---------------------------------------------------------------------------+
// | Child and Parent methods                                                  |
// +---------------------------------------------------------------------------+

func (o *keeper) Add(v Keeper) (yes bool) {
	o.mu.Lock()
	defer o.mu.Unlock()

	if v == nil {
		return
	}

	// 1. 重复添加.
	if _, ok := o.children[v.Name()]; ok {
		return false
	}

	// 2. 首次添加.
	v.SetParent(o)

	o.children[v.Name()] = v
	return true
}

func (o *keeper) Del(v Keeper) (yes bool) {
	o.mu.Lock()
	defer o.mu.Unlock()
	if _, ok := o.children[v.Name()]; ok {
		delete(o.children, v.Name())
		return true
	}
	return false
}

func (o *keeper) DelAuto(ad bool) Keeper { o.autoDelete = ad; return o }

func (o *keeper) SetParent(p Keeper) { o.parent = p }

// +---------------------------------------------------------------------------+
// | Event bound methods                                                       |
// +---------------------------------------------------------------------------+

func (o *keeper) After(handlers ...KeeperHandler) Keeper {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.listAfter = append(o.listAfter, handlers...)
	return o
}

func (o *keeper) Before(handlers ...KeeperHandler) Keeper {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.listBefore = append(o.listBefore, handlers...)
	return o
}

func (o *keeper) Listen(handlers ...KeeperHandler) Keeper {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.listListen = append(o.listListen, handlers...)
	return o
}

func (o *keeper) Panic(handler KeeperPanicHandler) Keeper {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.panicHandler = handler
	return o
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *keeper) childrenStart(ctx context.Context) {
	for _, v := range func() map[string]Keeper {
		o.mu.RLock()
		defer o.mu.RUnlock()
		return o.children
	}() {
		go func(k Keeper) {
			if err := k.Start(ctx); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("%s start child: %v", o.name, err))
			}
		}(v)
	}
}

func (o *keeper) childrenWaiter() bool {
	for _, v := range func() map[string]Keeper {
		o.mu.RLock()
		defer o.mu.RUnlock()
		return o.children
	}() {
		if !v.Stopped() {
			return o.childrenWaiter()
		}
	}
	return true
}

func (o *keeper) init() *keeper {
	o.children = make(map[string]Keeper)
	o.listAfter = make([]KeeperHandler, 0)
	o.listBefore = make([]KeeperHandler, 0)
	o.listListen = make([]KeeperHandler, 0)

	o.reset()
	return o
}

func (o *keeper) reset() {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.started = false
	o.restart = true
}

func (o *keeper) runHandlers(ctx context.Context, handlers []KeeperHandler) (ignored bool) {
	// 1. 捕获异常.
	defer func() {
		if v := recover(); v != nil {
			ignored = true

			// 触发执行器.
			if o.panicHandler != nil {
				o.panicHandler(ctx, v)
			} else {
				_, _ = fmt.Fprintf(os.Stderr, fmt.Sprintf("%s runtime fatal: %v", o.name, v))
			}
		}
	}()

	// 2. 遍历执行器.
	for _, handler := range handlers {
		if ignored = handler(ctx); ignored {
			break
		}
	}
	return
}
