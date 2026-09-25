package types

// TODO should this just be in internal/types?
type EventType string

const (
	EventTypeInstanceCreated EventType = "instance.created"
	EventTypeInstanceDeleted EventType = "instance.deleted"
)

type Event struct {
	Type EventType
	Data any
}

type EventInstanceCreated struct {
	ID string
}

type EventInstanceDeleted struct {
	ID       string
	Name     string
	Webhooks []string
}
