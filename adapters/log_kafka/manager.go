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

package log_kafka

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-wares/log/adapters"
	"github.com/go-wares/log/base"
	"github.com/go-wares/log/config"
	"sync"
	"time"
)

type (
	// Manager
	// 日志管理器.
	//
	// 发送用户日志到Kafka.
	Manager struct {
		bucket    *adapters.Bucket
		formatter adapters.LogFormatter
		keeper    base.Keeper
		mu        sync.RWMutex
		name      string
		producer  sarama.SyncProducer
	}
)

func New() adapters.LogAdapter {
	return (&Manager{}).init()
}

func (o *Manager) Keeper() base.Keeper { return o.keeper }

// Send
// 加入数据桶.
//
// 若数据桶积压数量超过指定值时, 立即刷盘保存.
func (o *Manager) Send(line *adapters.Line) {
	if n := o.bucket.Add(line); n >= config.Config.LogAdapterKafka.Batch {
		go o.save()
	}
}

// SetFormatter
// 设置格式.
func (o *Manager) SetFormatter(formatter adapters.LogFormatter) {
	o.formatter = formatter
}

// +---------------------------------------------------------------------------+
// | Event methods                                                             |
// +---------------------------------------------------------------------------+

func (o *Manager) onAfter(ctx context.Context) (ignored bool) {
	if o.bucket.Count() > 0 {
		o.save()
		return o.onAfter(ctx)
	}
	return
}

func (o *Manager) onListen(ctx context.Context) (ignored bool) {
	// 1. 定时保存.
	//    每隔指定时长(默认: 350ms)上报一次日志.
	ticker := time.NewTicker(time.Duration(config.Config.LogAdapterKafka.Milliseconds) * time.Millisecond)

	// 2. 关闭定时.
	defer ticker.Stop()

	// 3. 监听信号.
	for {
		select {
		case <-ticker.C:
			go o.save()
		case <-ctx.Done():
			return
		}
	}
}

// +---------------------------------------------------------------------------+
// | Access methods                                                            |
// +---------------------------------------------------------------------------+

func (o *Manager) init() *Manager {
	o.bucket = adapters.NewBucket()
	o.formatter = (&Formatter{}).init()
	o.name = fmt.Sprintf("log-kafka-manager")
	o.keeper = base.NewKeeper(o.name).
		After(o.onAfter).
		Listen(o.onListen)
	return o
}

func (o *Manager) getProducer() (sarama.SyncProducer, error) {
	o.mu.Lock()
	defer o.mu.Unlock()

	// 复用连接.
	if o.producer != nil {
		return o.producer, nil
	}

	// 准备连接.
	var (
		c   = sarama.NewConfig()
		err error
	)

	// 超时配置.
	c.Net.MaxOpenRequests = config.Config.LogAdapterKafka.ProducerMaxRequest
	c.Net.DialTimeout = time.Duration(config.Config.LogAdapterKafka.ProducerTimeout) * time.Second
	c.Net.ReadTimeout = time.Duration(config.Config.LogAdapterKafka.ProducerTimeout) * time.Second
	c.Net.WriteTimeout = time.Duration(config.Config.LogAdapterKafka.ProducerTimeout) * time.Second

	// 生产者配置.
	c.Producer.RequiredAcks = sarama.NoResponse
	c.Producer.Timeout = time.Duration(config.Config.LogAdapterKafka.ProducerTimeout) * time.Second
	c.Producer.Retry.Max = config.Config.LogAdapterKafka.ProducerRetry
	c.Producer.Retry.Backoff = 300 * time.Millisecond
	c.Producer.Return.Errors = true
	c.Producer.Return.Successes = true
	c.Producer.CompressionLevel = sarama.CompressionLevelDefault

	// 其它配置
	c.ChannelBufferSize = config.Config.LogAdapterKafka.ProducerBufferSize
	o.producer, err = sarama.NewSyncProducer(config.Config.LogAdapterKafka.Host, c)
	return o.producer, err
}

func (o *Manager) save() {
	var (
		list, count = o.bucket.Popn(config.Config.LogAdapterKafka.Batch)
		writer      *Writer
	)

	// 1. 空数据桶.
	if count == 0 {
		return
	}

	// 2. 释放实例.
	defer func() {
		// 2.1 释放日志.
		for _, v := range list {
			v.(*adapters.Line).Release()
		}

		// 2.2 释放实例.
		if writer != nil {
			writer.Release()
		}
	}()

	// 3. 获取实例.
	writer = NewWriter()
	writer.Send(o, list)
}
