package main

import (
	"github.com/ayeama/panel/api/internal/broker"
	"github.com/ayeama/panel/api/internal/scheduler"
)

// read agent stats
// listen for server create events
// move server create event from one stream to another based on agent stats

func main() {
	broker, err := broker.NewBroker(broker.BrokerTypeRedis)
	if err != nil {
		panic(err)
	}

	scheduler := scheduler.New(broker)
	scheduler.Start()
}
