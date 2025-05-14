package broker

import "fmt"

const (
	BrokerTypeRedis string = "redis"
)

type Broker interface{}

func NewBroker(t string) (Broker, error) {
	switch t {
	case BrokerTypeRedis:
		return NewBrokerRedis(), nil
	default:
		return nil, fmt.Errorf("unknown broker: %s", t)
	}
}
