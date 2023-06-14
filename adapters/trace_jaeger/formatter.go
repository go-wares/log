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
// date: 2023-02-24

package trace_jaeger

import (
	"context"
	"encoding/binary"
	"fmt"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/config"

	"github.com/go-wares/log/adapters/trace_jaeger/jaeger"
	"github.com/go-wares/log/adapters/trace_jaeger/thrift"
	"strconv"
)

type (
	formatter struct{}
)

func (o *formatter) Byte(vs ...adapters.Span) ([]byte, error) {
	return o.thrift(vs...)
}

// /////////////////////////////////////////////////////////////////////////////
// Access
// /////////////////////////////////////////////////////////////////////////////

func (o *formatter) build(list ...adapters.Span) (batch *jaeger.Batch) {
	return &jaeger.Batch{
		Process: o.buildProcess(),
		Spans:   o.buildSpans(list...),
	}
}

func (o *formatter) buildLogs(list []*adapters.Line) []*jaeger.Log {
	logs := make([]*jaeger.Log, 0)

	for _, x := range list {
		logs = append(logs, &jaeger.Log{
			Timestamp: x.Time.UnixMicro(),
			Fields: o.buildTagsMapper(x.Attr, adapters.Attr{
				x.Time.Format(config.Config.LogTimeFormat): x.Text,
				"log-level": x.Level,
			}),
		})
	}
	return logs
}

func (o *formatter) buildProcess() *jaeger.Process {
	return &jaeger.Process{
		ServiceName: config.Config.TraceAdapterJaeger.Topic,
		Tags:        o.buildTagsMapper(adapters.Resource),
	}
}

func (o *formatter) buildSpan(sp adapters.Span) *jaeger.Span {
	var (
		tid  = sp.Trace().TraceId()
		sid  = sp.SpanId()
		span = jaeger.NewSpan()
	)

	// Identify info.
	span.TraceIdHigh = int64(binary.BigEndian.Uint64(tid.Body()[0:8]))
	span.TraceIdLow = int64(binary.BigEndian.Uint64(tid.Body()[8:16]))
	span.SpanId = int64(binary.BigEndian.Uint64(sid.Body()[:]))

	if pid := sp.ParentSpanId(); pid != nil {
		span.ParentSpanId = int64(binary.BigEndian.Uint64(pid.Body()[:]))
	}

	// Basic info and flags.
	span.OperationName = sp.Name()
	span.StartTime = sp.StartTime().UnixMicro()
	span.Duration = sp.EndTime().Sub(sp.StartTime()).Microseconds()
	span.Flags = 1

	// Extensions.
	span.Tags = o.buildTagsMapper(sp.Attr())
	span.Logs = o.buildLogs(sp.Logs())
	span.References = o.buildReference()
	return span
}

func (o *formatter) buildSpans(sps ...adapters.Span) []*jaeger.Span {
	list := make([]*jaeger.Span, 0)
	for _, sp := range sps {
		list = append(list, o.buildSpan(sp))
	}
	return list
}

func (o *formatter) buildReference() (refs []*jaeger.SpanRef) { return nil }

func (o *formatter) buildTagsMapper(attrs ...adapters.Attr) []*jaeger.Tag {
	var (
		tags = make([]*jaeger.Tag, 0)
	)

	for _, attr := range attrs {
		for k, v := range attr {
			tag := &jaeger.Tag{Key: k}

			switch v.(type) {
			case bool:
				val := v.(bool)
				tag.VType = jaeger.TagType_BOOL
				tag.VBool = &val
			case float32, float64:
				val, _ := strconv.ParseFloat(fmt.Sprintf("%v", v), 64)
				tag.VType = jaeger.TagType_DOUBLE
				tag.VDouble = &val
			case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
				val, _ := strconv.ParseInt(fmt.Sprintf("%v", v), 10, 64)
				tag.VType = jaeger.TagType_LONG
				tag.VLong = &val
			case string:
				val := v.(string)
				tag.VType = jaeger.TagType_STRING
				tag.VStr = &val
			default:
				val := fmt.Sprintf("%v", v)
				tag.VType = jaeger.TagType_STRING
				tag.VStr = &val
			}

			tags = append(tags, tag)
		}
	}

	// Return
	// built tags.
	if len(tags) > 0 {
		return tags
	}
	return nil
}

func (o *formatter) init() *formatter { return o }

func (o *formatter) thrift(list ...adapters.Span) (buf []byte, err error) {
	var (
		bat = o.build(list...)
		ctx = context.Background()
		mem = thrift.NewTMemoryBuffer()
	)

	if err = bat.Write(ctx, thrift.NewTBinaryProtocolConf(mem, &thrift.TConfiguration{})); err == nil {
		buf = mem.Buffer.Bytes()
	}
	return
}
