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
// date: 2023-05-14

package tests

import (
	"github.com/go-wares/log/managers"
	"github.com/go-wares/log/trace"
	"strings"
	"testing"
	"time"
)

func TestSpan(t *testing.T) {
	defer managers.Manager.Stop()

	s10 := trace.NewSpan("span 10")
	defer s10.End()

	s10.Attr().Set("uid", 1).Set("key", "value")
	s10.Info("info message on trace")

	s20 := s10.Child("span 20")
	defer s20.End()

	s30 := s20.Child("span 30")
	defer s30.End()

	s40 := s30.Child("span 40")
	defer s40.End()

	s41 := s30.Child("span 41")
	defer s41.End()

}

func TestSpan_2(t *testing.T) {
	defer managers.Manager.Stop()

	sp := trace.NewSpan("start")
	t.Logf("sp: %32s | %16s | %16s", sp.Trace().TraceId(), strings.Repeat(" ", 16), sp.SpanId().String())
	time.Sleep(time.Millisecond * 10)
	defer sp.End()

	s1 := trace.NewSpanFromContext(sp.Context(), "span 1")
	t.Logf("s1: %32s | %16s | %16s", s1.Trace().TraceId(), s1.ParentSpanId().String(), s1.SpanId().String())
	time.Sleep(time.Millisecond * 10)
	defer s1.End()

	s2 := trace.NewSpanFromContext(s1.Context(), "span 2")
	t.Logf("s2: %32s | %16s | %16s", s2.Trace().TraceId(), s2.ParentSpanId().String(), s2.SpanId().String())
	time.Sleep(time.Millisecond * 10)
	defer s2.End()

	s3 := trace.NewSpanFromContext(s2.Context(), "span 3")
	t.Logf("s3: %32s | %16s | %16s", s3.Trace().TraceId(), s3.ParentSpanId().String(), s3.SpanId().String())
	time.Sleep(time.Millisecond * 10)
	defer s3.End()

	s4 := trace.NewSpanFromContext(s3.Context(), "span 4")
	t.Logf("s4: %32s | %16s | %16s", s4.Trace().TraceId(), s4.ParentSpanId().String(), s4.SpanId().String())
	time.Sleep(time.Millisecond * 10)
	s4.End()
}
