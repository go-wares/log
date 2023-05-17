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
// date: 2023-05-16

package tests

import (
	"github.com/Shopify/sarama"
	"testing"
	"time"
)

func TestKafka_Connect(t *testing.T) {
	var (
		cfg    = kafkaConfig()
		client sarama.Client
		err    error
		topic  string
		topics []string
	)

	// Build client.
	if client, err = sarama.NewClient(kafkaAddr(), cfg); err != nil {
		t.Errorf("[kafka client]: %v", err)
		return
	}

	// Broker list.
	for _, broker := range client.Brokers() {
		t.Logf("%v", broker)
	}

	// Topic list.
	if topics, err = client.Topics(); err != nil {
		t.Errorf("%v", err)
		return
	}

	// Topic range.
	t.Logf("[kafka topics]: %d", len(topics))
	for _, topic = range topics {
		t.Logf("[kafka  topic]: %v", topic)
	}
}

func TestKafka_SyncProducer(t *testing.T) {
	var (
		cfg      = kafkaConfig()
		client   sarama.Client
		err      error
		producer sarama.SyncProducer
	)

	cfg.Net.DialTimeout = time.Second
	cfg.Net.ReadTimeout = time.Second
	cfg.Net.WriteTimeout = time.Second

	cfg.Producer.Return.Errors = true
	cfg.Producer.Return.Successes = true

	if client, err = sarama.NewClient(kafkaAddr(), cfg); err != nil {
		t.Errorf("%v", err)
		return
	}

	t.Logf("client: %+v", client.Brokers()[0])

	if producer, err = sarama.NewSyncProducerFromClient(client); err != nil {
		t.Errorf("%v", err)
		return
	}

	var (
		partition int32
		offset    int64
	)

	if partition, offset, err = producer.SendMessage(&sarama.ProducerMessage{
		Topic: "go-wares-log",
		Value: sarama.StringEncoder(`{"service_addr": ["172.20.0.10", "192.168.0.130"]}`),
	}); err != nil {
		t.Errorf("%v", err)
		return
	}

	t.Logf("producer: partition=%d, offset=%d", partition, offset)
}

func kafkaAddr() []string {
	return []string{
		"192.168.0.130:9092",
		// "127.0.0.1:9092",
	}
}

func kafkaConfig() (cfg *sarama.Config) {
	cfg = sarama.NewConfig()
	return
}
