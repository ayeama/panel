package runtime

import (
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"strconv"
	"strings"

	"github.com/ayeama/panel/api/internal/config"
	"github.com/ayeama/panel/api/internal/domain"
)

type RuntimeType string

const RuntimeTypePodman RuntimeType = "podman"
const RuntimeTypeDocker RuntimeType = "docker"

type Runtime interface {
	New() (Runtime, error)

	Name(id string) string
	Status(id string) string
	Port(id string) string
	Running(id string) bool

	Create(id string, image string) string
	Delete(container *domain.Container)

	Start(container *domain.Container)
	Stop(container *domain.Container)
	Attach(container *domain.Container, stdin io.Reader, stdout io.Writer, stderr io.Writer, done chan struct{}) error
	Events() chan domain.RuntimeEvent
}

func New(t RuntimeType) (Runtime, error) {
	switch t {
	case RuntimeTypePodman:
		return Podman{}.New()
	case RuntimeTypeDocker:
		return Docker{}.New()
	default:
		return nil, fmt.Errorf("unknown runtime: %s", t)
	}
}

// TODO find better spot for this?
func freeHostPort() (uint16, error) {
	retries := 3
	var port uint16

	portRange := strings.Split(config.Config.ServerPortRange, "-")
	min, err := strconv.Atoi(portRange[0])
	if err != nil {
		panic(err)
	}
	max, err := strconv.Atoi(portRange[1])
	if err != nil {
		panic(err)
	}

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
