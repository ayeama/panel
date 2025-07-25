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

	Inspect(container_id string) domain.Container
	Create(id string, tag string) string
	Delete(container_id string)
	Start(container_id string)
	Stop(container_id string)
	Attach(container_id string, stdin io.Reader, stdout io.Writer, stderr io.Writer, done chan struct{}) error
	Events() chan domain.RuntimeEvent

	CreateSidecar(id string, tag string, server_id string) string
	InjectCredentials(container_id string)
}

func New(t RuntimeType) (Runtime, error) {
	switch t {
	case RuntimeTypePodman:
		return Podman{}.New()
	// case RuntimeTypeDocker:
	// 	return Docker{}.New()
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
