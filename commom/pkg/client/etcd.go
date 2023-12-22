package client

import (
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"google.dev/google/common/pkg/conf"
)

func ETCDClient(conf conf.ETCDConfiguration) (*clientv3.Client, error) {
	return clientv3.New(ETCDOption(conf))
}

func ETCDOption(conf conf.ETCDConfiguration) clientv3.Config {
	return clientv3.Config{
		Endpoints:   conf.Endpoints,
		DialTimeout: time.Duration(conf.DialTimeout) * time.Second,
		Username:    conf.Username,
		Password:    conf.Password,
	}
}
