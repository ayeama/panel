package types

type Type string

const (
	EventTypeInstance Type = "instance"
)

type Action string

const (
	EventActionCreate Action = "create"
	EventActionDelete Action = "delete"
)

type Actor struct {
	ID         string
	Attributes map[string]string
}

type Event struct {
	Type   Type
	Action Action
	Actor  Actor
}
