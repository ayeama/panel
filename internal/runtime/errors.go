package runtime

import (
	"errors"
	"fmt"
)

var (
	ErrInvalid  = errors.New("invalid")
	ErrNotFound = errors.New("not found")
	ErrExists   = errors.New("Exists")
)

type Error struct {
	Op       string
	Resource string
	ID       string
	Err      error
}

func (e Error) Error() string {
	return fmt.Sprintf(
		"%s %s %s: %v",
		e.Op,
		e.Resource,
		e.ID,
		e.Err,
	)
}

func (e Error) Unwrap() error {
	return e.Err
}
