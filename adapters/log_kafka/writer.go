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

package log_kafka

import (
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/config"
	"os"
	"sync"
)

var (
	writerPool sync.Pool
)

type (
	// Writer
	// 写日志.
	//
	// 上报日志到 Kafka.
	Writer struct{}
)

// NewWriter
// 获取写实例.
func NewWriter() *Writer {
	// 1. 池中获取.
	if g := writerPool.Get(); g != nil {
		return g.(*Writer).before()
	}

	// 2. 新建实例.
	g := (&Writer{}).init()
	return g.before()
}

// Release
// 释放实例.
func (o *Writer) Release() {
	o.after()
	writerPool.Put(o)
}

// Send
// 批量发送过程..
func (o *Writer) Send(manager *Manager, list []interface{}) {
	var (
		err      error
		msg      = make([]*sarama.ProducerMessage, 0)
		producer sarama.SyncProducer
		buf      []byte
	)

	// 1. 后置执行.
	defer func() {
		// 1.1 捕获异常.
		if v := recover(); v != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v, topic: %s, host: %v\n%s\n",
				v,
				config.Config.LogAdapterKafka.Topic,
				config.Config.LogAdapterKafka.Host,
				adapters.Backstack().String(),
			)
		}

		// 1.2 打印错误.
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "%v topic: %s, host: %v\n",
				err,
				config.Config.LogAdapterKafka.Topic,
				config.Config.LogAdapterKafka.Host,
			)
		}

		// 1.3 关闭连接.
		if producer != nil {
			if err = producer.Close(); err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v topic: %s, host: %v",
					err,
					config.Config.LogAdapterKafka.Topic,
					config.Config.LogAdapterKafka.Host,
				)
			}
		}
	}()

	// 2. 格式消息.
	for _, x := range list {
		if line, ok := x.(*adapters.Line); ok {
			// 2.1 消息正文.
			if buf = manager.formatter.Byte(line); buf == nil {
				continue
			}

			// 2.2 消息结构.
			msg = append(msg, &sarama.ProducerMessage{
				Topic: config.Config.LogAdapterKafka.Topic,
				Value: sarama.ByteEncoder(buf),
				// Key:       nil,
				// Headers:   nil,
				// Metadata:  nil,
				// Offset:    0,
				Partition: 0,
				// Timestamp: time.Time{},
			})
		}
	}

	// 3. 发送过程.
	if producer, err = manager.getProducer(); err == nil {
		err = producer.SendMessages(msg)
	}
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *Writer) after() *Writer  { return o }
func (o *Writer) before() *Writer { return o }
func (o *Writer) init() *Writer   { return o }
