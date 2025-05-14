package broker

import (
	"fmt"

	"github.com/ayeama/panel/api/internal/domain"
)

const (
	BrokerTypeRedis string = "redis"
)

type Broker interface {
	AddEventServerCreate(domain.EventServerCreate)
	AddEventServerDelete(domain.EventServerDelete)
	AddEventServerStart(domain.EventServerStart)
	AddEventServerStop(domain.EventServerStop)
}

func NewBroker(t string) (Broker, error) {
	switch t {
	case BrokerTypeRedis:
		return NewBrokerRedis(), nil
	default:
		return nil, fmt.Errorf("unknown broker: %s", t)
	}
}
