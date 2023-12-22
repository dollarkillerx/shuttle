package client

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"google.dev/google/common/pkg/conf"
)

// RedisClient ...
func RedisClient(conf conf.RedisConfiguration) (*redis.Client, error) {
	client := redis.NewClient(RedisOption(conf))

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*3)
	defer cancelFunc()
	err := client.Ping(ctx).Err()
	return client, err
}

// RedisOption ...
func RedisOption(conf conf.RedisConfiguration) *redis.Options {
	return &redis.Options{
		Addr:     conf.Addr,
		Password: conf.Password,
		DB:       conf.Db,

		//连接池容量及闲置连接数量
		//PoolSize:     15, // 连接池最大socket连接数，默认为4倍CPU数， Default is 10 connections per every available CPU as reported by runtime.GOMAXPROCS.
		MinIdleConns: 10, //在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量；。

		//超时
		DialTimeout:  5 * time.Second, //连接建立超时时间，默认5秒。
		ReadTimeout:  3 * time.Second, //读超时，默认3秒， -1表示取消读超时
		WriteTimeout: 3 * time.Second, //写超时，默认等于读超时
		PoolTimeout:  4 * time.Second, //当所有连接都处在繁忙状态时，客户端等待可用连接的最大等待时长，默认为读超时+1秒。

		//命令执行失败时的重试策略
		MaxRetries:      0,                      // 命令执行失败时，最多重试多少次，默认为0即不重试
		MinRetryBackoff: 8 * time.Millisecond,   //每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
		MaxRetryBackoff: 512 * time.Millisecond, //每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
	}
}
