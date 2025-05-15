package main

import (
	"github.com/ayeama/panel/api/internal/agent"
	"github.com/ayeama/panel/api/internal/broker"
	"github.com/ayeama/panel/api/internal/runtime"
)

func main() {
	broker, err := broker.NewBroker(broker.BrokerTypeRedis)
	if err != nil {
		panic(err)
	}

	runtime, err := runtime.NewRuntime(runtime.RuntimeTypePodman)
	if err != nil {
		panic(err)
	}

	agent := agent.New(broker, runtime)
	agent.Start()
}
