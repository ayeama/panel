package runtime

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net"

	"github.com/ayeama/panel/api/internal/domain"
)

const (
	RuntimeTypePodman string = "podman"
	RuntimeTypeDocker string = "docker"
)

type Runtime interface {
	Create(server *domain.Server) string
	Delete(container *domain.Container)
	Start(container *domain.Container)
	Stop(container *domain.Container)
	Stats(container *domain.Container) chan domain.ContainerStat
	Attach(container *domain.Container, stdin io.Reader, stdout io.Writer, stderr io.Writer) error
	Events() chan domain.Event
	Status(container *domain.Container) string
}

func New(t string) (Runtime, error) {
	switch t {
	case RuntimeTypePodman:
		return NewRuntimePodman(), nil
	default:
		return nil, fmt.Errorf("unknown runtim: %s", t)
	}
}

// TODO find better spot for this?
func freeHostPort() (uint16, error) {
	retries := 3
	var port uint16

	min := 10000
	max := 49152

	for i := 0; i < retries; i++ {
		port := uint16(rand.IntN(max-min) + min)
		listen, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
		if err != nil {
			continue
		}
		defer listen.Close()
		return port, err
	}

	return port, fmt.Errorf("panel: could not get free port in %d retries", retries)
}
