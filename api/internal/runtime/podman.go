package runtime

import (
	"context"
	"io"

	"github.com/ayeama/panel/api/internal/domain"
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

	resp, err := containers.CreateWithSpec(r.context, spec, nil)
	if err != nil {
		panic(err)
	}

	return resp.ID
}

func (r *Podman) Delete(container *domain.Container) {
	force := true
	volumes := true

	options := &containers.RemoveOptions{
		Force:   &force,
		Volumes: &volumes,
	}

	_, err := containers.Remove(r.context, container.Id, options)
	if err != nil {
		panic(err)
	}
}

func (r *Podman) Start(container *domain.Container) {
	err := containers.Start(r.context, container.Id, nil)
	if err != nil {
		panic(err)
	}
}

func (r *Podman) Stop(container *domain.Container) {
	err := containers.Stop(r.context, container.Id, nil)
	if err != nil {
		panic(err)
	}
}

func (r *Podman) Stats(container *domain.Container) chan domain.ContainerStat {
	all := false
	stream := true
	interval := 1

	options := &containers.StatsOptions{
		All:      &all,
		Stream:   &stream,
		Interval: &interval,
	}

	resp, err := containers.Stats(r.context, []string{container.Id}, options)
	if err != nil {
		panic(err)
	}

	stats := make(chan domain.ContainerStat)

	go func() {
		defer close(stats)

		for report := range resp {
			for _, stat := range report.Stats {
				stats <- domain.ContainerStat{
					Cpu:    stat.CPU,
					Memory: stat.MemPerc,
				}
			}
		}
	}()

	return stats
}

func (r *Podman) Attach(container *domain.Container, stdin io.Reader, stdout io.Writer, stderr io.Writer) {
	ready := make(chan bool)

	// TODO does this leak?
	go func() {
		err := containers.Attach(r.context, container.Id, stdin, stdout, stderr, ready, nil)
		if err != nil {
			panic(err)
		}
	}()

	<-ready
}

func (r *Podman) Status(container *domain.Container) string {
	resp, err := containers.Inspect(r.context, container.Id, nil)
	if err != nil {
		panic(err)
	}

	switch resp.State.Status {
	case "created":
		return "created"
	case "initialized":
		return "created"
	case "exited":
		return "stopped"
	case "paused":
		return "stopped"
	case "running":
		return "running"
	default:
		return "unknown"
	}
}
