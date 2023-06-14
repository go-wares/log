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
// date: 2023-05-17

package main

import (
	"encoding/json"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/go-wares/log/adapters/log_kafka"
	"time"
)

func main() {
	var (
		addr      = []string{"172.20.0.180:9092"}
		err       error
		cfg       = sarama.NewConfig()
		cli       sarama.Client
		offset    int64
		partition int32
		producer  sarama.SyncProducer
	)

	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true

	// Build client.
	if cli, err = sarama.NewClient(addr, cfg); err != nil {
		println("kafka build client: ", err.Error())
		return
	}

	// Build producer.
	if producer, err = sarama.NewSyncProducerFromClient(cli); err != nil {
		println("kafka build producer: %v", err.Error())
		return
	}

	curr := time.Now()
	line := &log_kafka.Data{
		Content: "log content",
		Keywords: map[string]interface{}{
			"id":  1,
			"key": "value",
		},
		Level:          "INFO",
		Time:           curr.Format("2006-01-02T15:04:05.999999Z"),
		TimestampMs:    curr.UnixMilli(),
		SpanId:         fmt.Sprintf("span-id-%d", curr.Unix()),
		TraceId:        fmt.Sprintf("trace-id-%d", curr.Unix()),
		RequestMethod:  "POST",
		RequestUrl:     "/user/register",
		UserAgent:      "Chrome/1.2.3",
		Pid:            3721,
		ServiceAddr:    []string{"172.20.0.100", "172.20.0.200"},
		ServiceName:    "go-wares-log",
		ServiceVersion: "1.0",
	}

	buf, _ := json.Marshal(line)

	// Send message to kafka.
	if partition, offset, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "go-wares-log",
		Value: sarama.ByteEncoder(buf),
	}); err != nil {
		println("kafka send error: %v", err.Error())
		return
	}

	println("kafka send completed: ", "partition =", partition, "offset =", offset)
}
