package main

import (
	"github.com/ayeama/panel/api/internal/broker"
	"github.com/ayeama/panel/api/internal/database"
	"github.com/ayeama/panel/api/internal/ep"
)

func main() {
	database := database.NewDatabase()

	broker, err := broker.NewBroker(broker.BrokerTypeRedis)
	if err != nil {
		panic(err)
	}

	ep := ep.New(database, broker)
	ep.Start()
}
