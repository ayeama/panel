package internal

import (
	"context"
	"log/slog"

	"github.com/containers/podman/v5/pkg/bindings"
)

type Agent struct {
	podman context.Context
}

func NewAgent() *Agent {
	return &Agent{}
}

func (a *Agent) Start() {
	slog.Info("starting")

	uri := "unix:///run/user/1000/podman/podman.sock"
	podman, err := bindings.NewConnection(context.Background(), uri)
	if err != nil {
		panic(err)
	}

	a.podman = podman
}
