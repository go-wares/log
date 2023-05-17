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

	var (
		err error
	)

	if o.producer != nil {
		return o.producer, nil
	}

	c := sarama.NewConfig()

	// c.Admin.Retry.Max = 5
	// c.Admin.Retry.Backoff = 100 * time.Millisecond
	// c.Admin.Timeout = 3 * time.Second

	// c.Net.MaxOpenRequests = 5
	// c.Net.DialTimeout = 30 * time.Second
	// c.Net.ReadTimeout = 30 * time.Second
	// c.Net.WriteTimeout = 30 * time.Second

	// c.Metadata.Retry.Max = 3
	// c.Metadata.Retry.Backoff = 250 * time.Millisecond
	// c.Metadata.RefreshFrequency = 10 * time.Minute
	// c.Metadata.Full = true
	// c.Metadata.AllowAutoTopicCreation = true

	// c.Producer.MaxMessageBytes = 1000000
	// c.Producer.RequiredAcks = sarama.WaitForLocal
	// c.Producer.Timeout = 10 * time.Second
	// c.Producer.Retry.Max = 3
	// c.Producer.Retry.Backoff = 100 * time.Millisecond
	// c.Producer.Return.Errors = true
	c.Producer.Return.Successes = true
	// c.Producer.CompressionLevel = sarama.CompressionLevelDefault
	// c.Producer.Transaction.Timeout = 1 * time.Minute
	// c.Producer.Transaction.Retry.Max = 50
	// c.Producer.Transaction.Retry.Backoff = 100 * time.Millisecond

	c.ChannelBufferSize = 256
	// c.ApiVersionsRequest = true
	// c.Version = sarama.DefaultVersion

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
