package client

import (
	"encoding/json"
	"fmt"
	"github.com/nsqio/go-nsq"
	"testing"
	"time"

	"google.dev/google/common/pkg/conf"
)

func TestNSQProducerClient(t *testing.T) {
	var cf = conf.NSQConfiguration{Address: []string{"127.0.0.1:4150"}}

	client, err := NSQProducerClient(cf)
	if err != nil {
		panic(err)
	}

	i := 0
	for {
		i++
		client.Publish("topic_test", []byte(fmt.Sprintf("hello world %d", i)))
		time.Sleep(time.Second)
	}
}

type messageHandler struct {
}

func (m *messageHandler) HandleMessage(message *nsq.Message) error {
	marshal, err := json.MarshalIndent(message, "", " ")
	if err == nil {
		fmt.Println(string(marshal))
	}
	return nil
}

func TestNSQConsumerClient(t *testing.T) {
	var cf = conf.NSQConfiguration{Address: []string{"127.0.0.1:4150"}}
	client, err := NSQConsumerClient(cf, "topic_test", "channel1")
	if err != nil {
		panic(err)
	}

	client.AddHandler(&messageHandler{})

	err = client.ConnectToNSQDs(cf.Address)
	if err != nil {
		panic(err)
	}

	for {
		time.Sleep(time.Second)
	}
}
