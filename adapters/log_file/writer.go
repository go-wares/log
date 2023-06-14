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

package log_file

import (
	"fmt"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/config"
	"os"
	"strings"
	"sync"
)

var (
	writerPool sync.Pool
)

type (
	Writer struct {
		paths map[string]bool
	}
)

func NewWriter() *Writer {
	if g := writerPool.Get(); g != nil {
		return g.(*Writer).before()
	}

	g := (&Writer{}).init()
	return g.before()
}

func (o *Writer) Release() {
	o.after()
	writerPool.Put(o)
}

func (o *Writer) Send(manager *Manager, list []interface{}) {
	var (
		files = make(map[string][]string)
	)

	// 1. 遍历日志.
	for _, v := range list {
		line := v.(*adapters.Line)

		// 1.1 文件夹名.
		//
		// - ./logs/2023-05/2023-05-13.log
		name := fmt.Sprintf("%s/%s/%s.%s",
			config.Config.LogAdapterFile.Path,
			line.Time.Format(config.Config.LogAdapterFile.Folder),
			line.Time.Format(config.Config.LogAdapterFile.Name),
			config.Config.LogAdapterFile.Ext,
		)

		// 1.2 创建目录.
		if _, ok := files[name]; !ok {
			files[name] = make([]string, 0)
			manager.mkdir(fmt.Sprintf("%s/%s",
				config.Config.LogAdapterFile.Path,
				line.Time.Format(config.Config.LogAdapterFile.Folder),
			))
		}

		// 1.3 加入日志.
		files[name] = append(files[name], manager.formatter.String(line))
	}

	// 2. 并行写入.
	w := &sync.WaitGroup{}
	for fp, fl := range files {
		w.Add(1)
		go func(fp string, fl []string) {
			defer w.Done()
			o.write(fp, fl)
		}(fp, fl)
	}
	w.Wait()
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *Writer) after() *Writer {
	return o
}

func (o *Writer) before() *Writer {
	return o
}

func (o *Writer) init() *Writer {
	o.paths = make(map[string]bool)
	return o
}

// 写入日志.
func (o *Writer) write(path string, list []string) {
	var (
		err  error
		file *os.File
	)

	// 1. 打开文件.
	if file, err = os.OpenFile(path, os.O_RDWR|os.O_APPEND|os.O_CREATE, os.ModePerm); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "file open: %v\n", err)
		return
	}

	// 2. 关闭文件.
	defer func() {
		if err = file.Close(); err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "file close: %v\n", err)
		}
	}()

	// 3. 写入日志.
	if _, err = file.WriteString(strings.Join(list, "\n") + "\n"); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "file write: %v\n", err)
	}
}
