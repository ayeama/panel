package broker

import (
	"context"
	"encoding/json"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/redis/go-redis/v9"
)

type BrokerRedis struct {
	rdb *redis.Client
}

func NewBrokerRedis() *BrokerRedis {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	})
	return &BrokerRedis{rdb: rdb}
}

func (b *BrokerRedis) AddEventServerCreate(e domain.EventServerCreate) {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	_, err = b.rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "server:create",
		Values: map[string]interface{}{
			"data": data,
		},
	}).Result()
	if err != nil {
		panic(err)
	}
}

func (b *BrokerRedis) AddEventServerDelete(e domain.EventServerDelete) {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	_, err = b.rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "server:delete",
		Values: map[string]interface{}{
			"data": data,
		},
	}).Result()
	if err != nil {
		panic(err)
	}
}

func (b *BrokerRedis) AddEventServerStart(e domain.EventServerStart) {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	// todo change stram
	_, err = b.rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "server:start",
		Values: map[string]interface{}{
			"data": data,
		},
	}).Result()
	if err != nil {
		panic(err)
	}
}

func (b *BrokerRedis) AddEventServerStop(e domain.EventServerStop) {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	// todo change stram
	_, err = b.rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "server:stop",
		Values: map[string]interface{}{
			"data": data,
		},
	}).Result()
	if err != nil {
		panic(err)
	}
}
