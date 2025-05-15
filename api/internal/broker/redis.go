package broker

import (
	"context"
	"encoding/json"
	"fmt"

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

func (b *BrokerRedis) AddEventAgentCommand(e domain.EventServerCreate) {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	_, err = b.rdb.XAdd(context.Background(), &redis.XAddArgs{
		Stream: "agent:neon:commands", // TODO variablise
		Values: map[string]interface{}{
			"data": data,
		},
	}).Result()
	if err != nil {
		panic(err)
	}
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

func (b *BrokerRedis) PublishEventAgentStat(e domain.EventAgentStat) {
	data, err := json.Marshal(e)
	if err != nil {
		panic(err)
	}

	err = b.rdb.Publish(context.Background(), "agent:stat", data).Err()
	if err != nil {
		panic(err)
	}
}

func (b *BrokerRedis) ReadEventAgentCommand() domain.EventServerCreate {
	hostname := "neon" // TODO pass in agent id?

	// TODO update stream to agent specific?
	err := b.rdb.XGroupCreateMkStream(context.Background(), "agent:neon:commands", "agent", "0").Err()
	if err != nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			panic(err)
		}
	}

	r, err := b.rdb.XReadGroup(context.Background(), &redis.XReadGroupArgs{
		Group:    "agent",
		Consumer: hostname,
		Streams:  []string{"agent:neon:commands", ">"},
		Count:    1,
		Block:    0,
	}).Result()
	if err != nil {
		panic(err)
	}

	s := r[0]
	m := s.Messages[0]

	d, ok := m.Values["data"]
	if !ok {
		panic(fmt.Errorf("cound not find data in stream message")) // TODO better errors?
	}

	data, ok := d.(string)
	if !ok {
		panic(fmt.Errorf("expected string, got %T", d)) // TODO better errors?
	}

	// TODO handle multiple different events?
	var event domain.EventServerCreate
	err = json.Unmarshal([]byte(data), &event)
	if err != nil {
		panic(err)
	}

	return event
}

func (b *BrokerRedis) ReadEventAgentStat() domain.EventAgentStat {
	ps := b.rdb.Subscribe(context.Background(), "agent:stat")
	defer ps.Close() // TODO this is because there were >600 connected clients to redis :o
	m, err := ps.ReceiveMessage(context.Background())
	if err != nil {
		panic(err)
	}

	var event domain.EventAgentStat
	err = json.Unmarshal([]byte(m.Payload), &event)
	if err != nil {
		panic(err)
	}

	return event
}

func (b *BrokerRedis) ReadEventServerCreate() domain.EventServerCreate {
	err := b.rdb.XGroupCreateMkStream(context.Background(), "server:create", "scheduler", "0").Err()
	if err != nil {
		if err.Error() != "BUSYGROUP Consumer Group name already exists" {
			panic(err)
		}
	}

	r, err := b.rdb.XReadGroup(context.Background(), &redis.XReadGroupArgs{
		Group:    "scheduler",
		Consumer: "scheduler",
		Streams:  []string{"server:create", ">"},
		Count:    1,
		Block:    0,
	}).Result()
	if err != nil {
		panic(err)
	}

	s := r[0]
	m := s.Messages[0]

	d, ok := m.Values["data"]
	if !ok {
		panic(fmt.Errorf("cound not find data in stream message")) // TODO better errors?
	}

	data, ok := d.(string)
	if !ok {
		panic(fmt.Errorf("expected string, got %T", d)) // TODO better errors?
	}

	// TODO handle multiple different events?
	var event domain.EventServerCreate
	err = json.Unmarshal([]byte(data), &event)
	if err != nil {
		panic(err)
	}

	return event
}
