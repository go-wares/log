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

package adapters

import (
	"sync"
)

type (
	// Bucket
	// 日志数据桶.
	Bucket struct {
		lines []interface{}
		mu    sync.RWMutex
	}
)

// NewBucket
// 创建数据桶.
func NewBucket() *Bucket {
	return &Bucket{
		lines: make([]interface{}, 0),
	}
}

// Add
// 添加日志入桶.
func (o *Bucket) Add(lines ...interface{}) int {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.lines = append(o.lines, lines...)
	return len(o.lines)
}

// Count
// 获取桶日志数量.
func (o *Bucket) Count() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return len(o.lines)
}

// Pop
// 取出1条日志.
func (o *Bucket) Pop() interface{} {
	if v, count := o.Popn(1); count == 1 {
		return v[0]
	}
	return nil
}

// Popn
// 取出N条日志.
func (o *Bucket) Popn(n int) (list []interface{}, count int) {
	o.mu.Lock()
	defer o.mu.Unlock()

	// 1. 取出全部.
	if total := len(o.lines); n >= total {
		// 1.1 空数据桶.
		if total == 0 {
			return
		}

		// 1.2 全部取出.
		count = total
		list = o.lines[0:]

		// 1.3 重置数据桶.
		o.lines = make([]interface{}, 0)
		return
	}

	// 2. 取出片段.
	count = n
	list = o.lines[0:n]

	// 2.1 重置数据桶.
	o.lines = o.lines[n:]
	return
}
