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
// date: 2023-05-13

package tests

import (
	"context"
	"github.com/go-wares/log"
	"github.com/go-wares/log/base"
	"github.com/go-wares/log/managers"
	"testing"
	"time"
)

type keeper struct {
	name string
}

func (o *keeper) onAfter(ctx context.Context) (ignored bool) {
	log.Infof("%s - event:after", o.name)
	return
}

func (o *keeper) onBefore(ctx context.Context) (ignored bool) {
	log.Infof("%s - event:before", o.name)
	return
}

func (o *keeper) onListen(ctx context.Context) (ignored bool) {
	log.Infof("%s - event:listen", o.name)
	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}

func (o *keeper) onPanic(ctx context.Context, v interface{}) {
	log.Fatalf("%s - event:panic: %v", o.name, v)
	return
}

func TestKeeper(t *testing.T) {
	var (
		ctx, cancel = context.WithCancel(context.TODO())
		keeper1     = &keeper{name: "v1"}
		keeper2     = &keeper{name: "v1.1"}
		keeper3     = &keeper{name: "v1.2"}
		nk1         base.Keeper
	)

	defer managers.Manager.Stop()

	go func() {
		time.Sleep(time.Second * 3)
		nk1.Restart()

		time.Sleep(time.Second * 3)
		cancel()
	}()

	nk1 = base.NewKeeper(keeper1.name)
	nk1.After(keeper1.onAfter).Before(keeper1.onBefore).Listen(keeper1.onListen).Panic(keeper1.onPanic)

	nk2 := base.NewKeeper(keeper2.name)
	nk2.After(keeper2.onAfter).Before(keeper2.onBefore).Listen(keeper2.onListen).Panic(keeper2.onPanic)
	nk1.Add(nk2)

	nk3 := base.NewKeeper(keeper3.name)
	nk3.After(keeper3.onAfter).Before(keeper3.onBefore).Listen(keeper3.onListen).Panic(keeper3.onPanic)
	nk1.Add(nk3)

	nk1.Start(ctx)
}
