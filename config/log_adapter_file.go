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

package config

type (
	// LogAdapterFile
	// 文件适配器配置.
	//
	//   # config/log.yaml
	//
	//   log_adapter: file
	//   log_adapter_file:
	//     ext: log
	//     path: ./logs
	//     folder: 2006-01
	//     name: 2006-01-02.log
	LogAdapterFile struct {
		// 批量阈值.
		// 每次最多批量写入N(默认: 100)条日志.
		Batch int `yaml:"batch" json:"batch"`

		// 保时频率.
		// 每隔固定时长(默认: 350ms)刷盘一次日志.
		Milliseconds int64 `yaml:"milliseconds" json:"milliseconds"`

		// 存储位置.
		//
		// - 默认：./logs
		// - 说明：日志文件存储到哪个位置, 物理路径.
		Path string `yaml:"path" json:"path"`

		// 文件夹拆分.
		//
		// - 默认：2006-01 (即月份, 如: 2023-05)
		// - 说明：日志文件按时间拆分目录.
		Folder string `yaml:"folder" json:"folder"`

		// 日志文件名.
		//
		// - 默认：2006-01-02 (即日期, 如: 2023-05-13)
		Name string `yaml:"name" json:"name"`

		// 日志扩展名.
		//
		// - 默认：log
		Ext string `yaml:"ext" json:"ext"`
	}
)

func (o *LogAdapterFile) defaults(_ *Configuration) {
	if o.Batch == 0 {
		o.Batch = defaultLogAdapterFileBatch
	}
	if o.Milliseconds == 0 {
		o.Milliseconds = defaultLogAdapterFileMilliseconds
	}
	if o.Path == "" {
		o.Path = defaultLogAdapterFilePath
	}
	if o.Folder == "" {
		o.Folder = defaultLogAdapterFileFolder
	}
	if o.Name == "" {
		o.Name = defaultLogAdapterFileName
	}
	if o.Ext == "" {
		o.Ext = defaultLogAdapterFileExt
	}
}
