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
)

const (
	OpenTelemetrySpan  = "__OPEN_TELEMETRY_SPAN__"
	OpenTelemetryTrace = "__OPEN_TELEMETRY_TRACE__"

	OpenTracingKey          = "__OPEN_TRACING_KEY__"
	OpenTracingParentSpanId = "X-B3-Parentspanid"
	OpenTracingSpanId       = "X-B3-Spanid"
	OpenTracingTraceId      = "X-B3-Traceid"
	OpenTracingSampled      = "X-B3-Sampled"
	OpenTracingSampledFlag  = "1"
)

var (
	defaultAutoStart           = true
	defaultLogAdapterTermColor = true
)

const (
	defaultLogAdapter = base.Term

	defaultLogAdapterFileBatch        = 100
	defaultLogAdapterFileMilliseconds = 350
	defaultLogAdapterFileExt          = "log"
	defaultLogAdapterFilePath         = "./logs"
	defaultLogAdapterFileFolder       = "2006-01"
	defaultLogAdapterFileName         = "2006-01-02"

	defaultLogAdapterKafkaHost  = "127.0.0.1:9092"
	defaultLogAdapterKafkaTopic = "logs"

	defaultLogTimeFormat = "2006-01-02 15:04:05.999"

	defaultTraceAdapterJaegerBatch        = 100
	defaultTraceAdapterJaegerMilliseconds = 350
	defaultTraceAdapterJaegerTopic        = "logs"
)
