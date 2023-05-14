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
	"testing"
)

func TestSpan(t *testing.T) {
	defer managers.Manager.Stop()

	s1 := trace.NewSpan("parent")
	defer s1.End()

	s1.Attr().Set("uid", 1).Set("key", "value")
	s1.Info("info message on trace")

	s2 := s1.Child("child s2")
	defer s2.End()

	s3 := s1.Child("child s3")
	defer s3.End()

}
