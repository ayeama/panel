package runtime

import (
	"context"

	"github.com/containers/podman/v5/pkg/bindings"
)

type RuntimePodman struct {
	Context context.Context
}

func NewRuntimePodman() *RuntimePodman {
	uri := "unix:///run/user/1000/podman/podman.sock"
	context, err := bindings.NewConnection(context.Background(), uri)
	if err != nil {
		panic(err)
	}

	return &RuntimePodman{Context: context}
}
