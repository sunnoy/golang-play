/**
    * @Author lirui
    * @Date 2022/3/24 14:44
**/

package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

type SentinelConn struct {
	MasterName    string
	SentinelAddrs []string
	Password      string
	DB            int
	Client        *redis.Client
}

func (sentinel *SentinelConn) GetConn() (*redis.Client, error) {
	if sentinel.Client == nil {
		return nil, fmt.Errorf("redis Client connect is nil")
	}
	ctx := context.Background()
	err := sentinel.Client.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}
	return sentinel.Client, nil
}

func (sentinel *SentinelConn) InitSentinelRedis() *redis.Client {
	sf := &redis.FailoverOptions{
		MasterName:    sentinel.MasterName,
		SentinelAddrs: sentinel.SentinelAddrs,
		Password:      sentinel.Password,
		DB:            sentinel.DB,
	}

	return redis.NewFailoverClient(sf)
}

func NewRedisSentinel(mastername, password string, addrs []string, db int) *SentinelConn {
	sentinel := SentinelConn{
		MasterName:    mastername,
		SentinelAddrs: addrs,
		Password:      password,
		DB:            db,
	}
	sentinel.Client = sentinel.InitSentinelRedis()

	return &sentinel

}
