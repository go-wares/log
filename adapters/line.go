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

package adapters

import (
	"context"
	"fmt"
	"github.com/go-wares/log/base"
	"sync"
	"time"
)

var (
	linePool sync.Pool
)

type (
	// Line
	// 单行日志.
	Line struct {
		Attr  Attr
		Ctx   context.Context
		Level base.LogLevel
		Text  string
		Time  time.Time

		SpanId, ParentSpanId string
		TraceId              string
	}
)

func NewLine(ctx context.Context, level base.LogLevel, format string, args ...interface{}) *Line {
	if x := linePool.Get(); x != nil {
		return x.(*Line).before(ctx, level, format, args...)
	}

	x := (&Line{}).init()
	return x.before(ctx, level, format, args...)
}

func (o *Line) Release() {
	o.after()
	linePool.Put(o)
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *Line) after() *Line {
	o.Attr = nil
	o.Ctx = nil
	o.Level = base.Off
	o.Text = ""
	return o
}

func (o *Line) before(ctx context.Context, level base.LogLevel, format string, args ...interface{}) *Line {
	o.Ctx = ctx
	o.Level = level
	o.Time = time.Now()
	o.Text = fmt.Sprintf(format, args...)

	// 堆椎清单.
	if level == base.Fatal {
		o.Text = fmt.Sprintf("%s\n%s", o.Text, Backstack().InternalString())
	}

	if ctx != nil {
		o.openTracing(ctx)
	}

	return o
}

func (o *Line) init() *Line {
	return o
}

func (o *Line) openTracing(ctx context.Context) {

}
