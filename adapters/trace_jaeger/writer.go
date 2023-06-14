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

package trace_jaeger

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/config"
	"github.com/valyala/fasthttp"
	"net/http"
	"os"
	"sync"
)

var (
	writerPool sync.Pool
)

type (
	Writer interface {
		Release()
		Send(formatter *formatter, lines ...adapters.Span)
	}

	writer struct {
		request  *fasthttp.Request
		response *fasthttp.Response
	}
)

func NewWriter() Writer {
	if o := writerPool.Get(); o != nil {
		return o.(*writer).before()
	}

	o := (&writer{}).init()
	o.before()
	return o
}

func (o *writer) Release() {
	o.after()
	writerPool.Put(o)
}

// Send
// 发送链路消息.
func (o *writer) Send(formatter *formatter, lines ...adapters.Span) {
	// 1. 后置执行.
	defer func() {
		// 1.1 捕获异常.
		if r := recover(); r != nil {
			_, _ = fmt.Fprintf(os.Stderr, "jaeger fatal: %v\n%s\n", r,
				adapters.Backstack().String(),
			)
		}
	}()

	// 2. 格式转换.
	var (
		body []byte
		err  error
	)
	if body, err = formatter.Byte(lines...); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "jaeger formatter: %v\n", err)
	}

	// 2.1 构建消息.
	buf := bytes.NewBuffer(body)

	// 3. 准备请求.
	o.request.SetRequestURI(config.Config.TraceAdapterJaeger.Endpoint)
	o.request.SetBodyStream(buf, buf.Len())
	o.request.Header.SetMethod(http.MethodPost)
	o.request.Header.SetContentType("application/x-thrift")

	// 4. 基础鉴权.
	if usr := config.Config.TraceAdapterJaeger.Username; usr != "" {
		pwd := config.Config.TraceAdapterJaeger.Password
		o.request.Header.Set("Authorization", fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString([]byte(usr+":"+pwd))))
	}

	// 5. 发送请求.
	if err = fasthttp.Do(o.request, o.response); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "jaeger trace: %v\n", err)
	}
	return
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *writer) after() *writer {
	fasthttp.ReleaseRequest(o.request)
	fasthttp.ReleaseResponse(o.response)

	o.request = nil
	o.response = nil
	return o
}

func (o *writer) before() *writer {
	o.request = fasthttp.AcquireRequest()
	o.response = fasthttp.AcquireResponse()
	return o
}

func (o *writer) init() *writer {
	return o
}
