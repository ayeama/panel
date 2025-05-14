package internal

import (
	"log/slog"

	"github.com/ayeama/panel/agent/internal/broker"
	"github.com/ayeama/panel/agent/internal/runtime"
)

type Agent struct {
	Broker  *broker.Broker
	Runtime *runtime.Runtime
}

func NewAgent() *Agent {
	broker, err := broker.NewBroker(broker.BrokerTypeRedis)
	if err != nil {
		panic(err)
	}

	runtime, err := runtime.NewRuntime(runtime.RuntimeTypePodman)
	if err != nil {
		panic(err)
	}

	return &Agent{Broker: &broker, Runtime: &runtime}
}

func (a *Agent) Start() {
	slog.Info("starting")
}
