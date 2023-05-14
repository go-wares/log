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

package config

import (
	"github.com/go-wares/log/base"
	"gopkg.in/yaml.v3"
	"os"
)

var (
	// Config
	// 配置实例.
	Config *Configuration
)

type (
	// Configuration
	// 基础配置.
	Configuration struct {
		// 自动启动.
		// 当包加载完成后, 立即启动监听.
		AutoStart *bool `yaml:"auto_start" json:"auto_start"`

		// 日志级别.
		//
		// - 默认: INFO
		// - 支持: DEBUG, INFO, WARN, ERROR, FATAL, OFF
		Level         base.Level    `yaml:"level" json:"level"`
		LogLevel      base.LogLevel `yaml:"-" json:"log_level"`
		LogTimeFormat string        `yaml:"log_time_format" json:"log_time_format"`

		// 日志适配器.
		//
		// - 默认：term
		// - 支持：term, file, kafka
		LogAdapter                                base.LogAdapter  `yaml:"log_adapter" json:"log_adapter"`
		LogAdapterTerm                            *LogAdapterTerm  `yaml:"log_adapter_term" json:"log_adapter_term"`
		LogAdapterFile                            *LogAdapterFile  `yaml:"log_adapter_file" json:"log_adapter_file"`
		LogAdapterKafka                           *LogAdapterKafka `yaml:"log_adapter_kafka" json:"log_adapter_kafka"`
		debugOn, infoOn, warnOn, errorOn, fatalOn bool

		// 链路适配器.
		//
		// - 默认：无
		// - 支持：jaeger, zipkin
		TraceAdapter       base.TraceAdapter   `yaml:"trace_adapter" json:"trace_adapter"`
		TraceAdapterJaeger *TraceAdapterJaeger `yaml:"trace_adapter_jaeger" json:"trace_adapter_jaeger"`
		TraceAdapterZipkin *TraceAdapterZipkin `yaml:"trace_adapter_zipkin" json:"trace_adapter_zipkin"`
	}
)

// +---------------------------------------------------------------------------+
// | Switch methods                                                            |
// +---------------------------------------------------------------------------+

func (o *Configuration) DebugOn() bool { return o.debugOn }
func (o *Configuration) InfoOn() bool  { return o.infoOn }
func (o *Configuration) WarnOn() bool  { return o.warnOn }
func (o *Configuration) ErrorOn() bool { return o.errorOn }
func (o *Configuration) FatalOn() bool { return o.fatalOn }

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *Configuration) defaults() {
	if o.AutoStart == nil {
		o.AutoStart = &defaultAutoStart
	}

	// 日志级别.
	o.Level, o.LogLevel = o.Level.LogLevel()
	o.debugOn = o.LogLevel >= base.Debug // 4 >= 5
	o.infoOn = o.LogLevel >= base.Info   // 4 >= 4
	o.warnOn = o.LogLevel >= base.Warn   // 4 >= 3
	o.errorOn = o.LogLevel >= base.Error // 4 >= 2
	o.fatalOn = o.LogLevel >= base.Fatal // 4 >= 1

	// 时间格式.
	if o.LogTimeFormat == "" {
		o.LogTimeFormat = defaultLogTimeFormat
	}

	// 适配器.
	if o.LogAdapter == "" {
		o.LogAdapter = defaultLogAdapter
	}

	// 终端适配器.
	// 当 LogAdapter 值为 term 时以下配置生效.
	if o.LogAdapterTerm == nil {
		o.LogAdapterTerm = &LogAdapterTerm{}
	}
	o.LogAdapterTerm.defaults()

	// 文件适配器.
	if o.LogAdapterFile == nil {
		o.LogAdapterFile = &LogAdapterFile{}
	}
	o.LogAdapterFile.defaults()

	// 消息适配器/Kafka.
	if o.LogAdapterKafka == nil {
		o.LogAdapterKafka = &LogAdapterKafka{}
	}
	o.LogAdapterKafka.defaults()

	// Jaeger 适配器.
	if o.TraceAdapterJaeger == nil {
		o.TraceAdapterJaeger = &TraceAdapterJaeger{}
	}
	o.TraceAdapterJaeger.defaults()

	// Zipkin 适配器.
	if o.TraceAdapterZipkin == nil {
		o.TraceAdapterZipkin = &TraceAdapterZipkin{}
	}
	o.TraceAdapterZipkin.defaults()
}

func (o *Configuration) init() *Configuration {
	o.scan().defaults()
	return o
}

func (o *Configuration) scan() *Configuration {
	for _, path := range []string{"config/log.yaml", "../config/log.yaml"} {
		body, err := os.ReadFile(path)
		if err == nil {
			if yaml.Unmarshal(body, o) == nil {
				break
			}
		}
	}
	return o
}
