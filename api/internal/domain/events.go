package domain

type Event interface{}

type ServerCreatedEvent struct {
	Id          string
	ContainerId string
}

type ServerStartedEvent struct {
	Id          string
	ContainerId string
}

type ServerStoppedEvent struct {
	Id          string
	ContainerId string
}

type ServerDeletedEvent struct {
	Id          string
	ContainerId string
}
