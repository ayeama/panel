package runtime

import "fmt"

const (
	RuntimeTypePodman string = "podman"
	RuntimeTypeDocker string = "docker"
)

type Runtime interface{}

func NewRuntime(t string) (Runtime, error) {
	switch t {
	case RuntimeTypePodman:
		return NewRuntimePodman(), nil
	default:
		return nil, fmt.Errorf("unknown runtime: %s", t)
	}
}
