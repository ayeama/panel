package broker

import "github.com/redis/go-redis/v9"

type BrokerRedis struct {
	Rdb *redis.Client
}

func NewBrokerRedis() *BrokerRedis {
	options := redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}
	rdb := redis.NewClient(&options)

	return &BrokerRedis{Rdb: rdb}
}
