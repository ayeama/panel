package runtime

import (
	"context"

	nettypes "github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type Podman struct {
	context context.Context
}

func NewRuntimePodman() *Podman {
	uri := "unix:///run/user/1000/podman/podman.sock"
	context, err := bindings.NewConnection(context.Background(), uri)
	if err != nil {
		panic(err)
	}

	return &Podman{context: context}
}

func (r *Podman) Create() string {
	stdin := true
	terminal := true

	cpus := 1.0
	cpuPeriod := uint64(100000)
	cpuQuota := int64(float64(cpuPeriod) * cpus)
	memLimit := int64(1000000000)

	hostPort, err := freeHostPort()
	if err != nil {
		panic(err)
	}

	var portMappings []nettypes.PortMapping
	portMappings = append(portMappings, nettypes.PortMapping{HostPort: hostPort, ContainerPort: 25565})

	spec := specgen.NewSpecGenerator("localhost/server/minecraft:latest", false)
	spec.Stdin = &stdin
	spec.Terminal = &terminal
	spec.ResourceLimits = &specs.LinuxResources{
		CPU: &specs.LinuxCPU{
			Period: &cpuPeriod,
			Quota:  &cpuQuota,
		},
		Memory: &specs.LinuxMemory{
			Limit: &memLimit,
		},
	}
	spec.PortMappings = portMappings

	// TODO add metadata to labels
	// spec.Labels = map[string]string{"com.github.ayeama.panel.api.server.id": "id"}

	// todo container restart: unless-stopped

	createResponse, err := containers.CreateWithSpec(r.context, spec, nil)
	if err != nil {
		panic(err)
	}

	return createResponse.ID
}

func (r *Podman) Delete() {}

func (r *Podman) Start() {}

func (r *Podman) Stop() {}
