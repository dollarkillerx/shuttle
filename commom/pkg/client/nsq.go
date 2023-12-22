package client

import (
	"github.com/nsqio/go-nsq"
	"google.dev/google/common/pkg/conf"
)

func NSQProducerClient(conf conf.NSQConfiguration) (*nsq.Producer, error) {
	config := nsq.NewConfig()
	return nsq.NewProducer(conf.Address[0], config)
}

func NSQConsumerClient(conf conf.NSQConfiguration, topic string, channel string) (*nsq.Consumer, error) {
	config := nsq.NewConfig()
	return nsq.NewConsumer(topic, channel, config)
}
