package runtime

import (
	"context"
	"fmt"
	"io"
	"math/rand/v2"
	"net"
	"strconv"
	"strings"

	"github.com/ayeama/panel/api/internal/config"
	"github.com/ayeama/panel/api/internal/domain"
	nettypes "github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/libpod/define"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/system"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/containers/podman/v5/pkg/specgen"
	dockerevents "github.com/docker/docker/api/types/events"
)

type Runtime struct {
	ctx context.Context
}

func New() (*Runtime, error) {
	ctx, err := bindings.NewConnection(context.Background(), config.Config.RuntimeUri)
	if err != nil {
		return nil, err
	}
	return &Runtime{ctx: ctx}, nil
}

func (r *Runtime) inspect(id string) *define.InspectContainerData {
	resp, err := containers.Inspect(r.ctx, id, nil)
	if err != nil {
		panic(err)
	}
	return resp
}

func (r *Runtime) Name(id string) string {
	resp := r.inspect(id)
	return resp.Name
}

func (r *Runtime) Status(id string) string {
	resp := r.inspect(id)
	return resp.State.Status
}

func (r *Runtime) Port(id string) string {
	resp := r.inspect(id)
	for _, bindings := range resp.HostConfig.PortBindings {
		for _, binding := range bindings {
			return binding.HostPort
		}
	}
	return ""
}

func (r *Runtime) Running(id string) bool {
	resp := r.inspect(id)
	return resp.State.Running
}

func (r *Runtime) Create(id string, image string) string {
	stdin := true
	terminal := true

	// cpus := 1.0
	// cpuPeriod := uint64(100000)
	// cpuQuota := int64(float64(cpuPeriod) * cpus)
	// memLimit := int64(1000000000)

	hostPort, err := freeHostPort()
	if err != nil {
		panic(err)
	}

	var portMappings []nettypes.PortMapping
	portMappings = append(portMappings, nettypes.PortMapping{HostPort: hostPort, ContainerPort: 25565})

	spec := specgen.NewSpecGenerator(image, false)
	spec.Stdin = &stdin
	spec.Terminal = &terminal
	// spec.ResourceLimits = &specs.LinuxResources{
	// 	CPU: &specs.LinuxCPU{
	// 		Period: &cpuPeriod,
	// 		Quota:  &cpuQuota,
	// 	},
	// 	Memory: &specs.LinuxMemory{
	// 		Limit: &memLimit,
	// 	},
	// }
	spec.PortMappings = portMappings

	// restartRetries := uint(3)
	// spec.RestartPolicy = "unless-stopped"
	// spec.RestartRetries = &restartRetries

	spec.Labels = make(map[string]string)
	spec.Labels["com.github.ayeama.panel.api.server.id"] = id

	resp, err := containers.CreateWithSpec(r.ctx, spec, nil)
	if err != nil {
		panic(err)
	}

	return resp.ID
}

func (r *Runtime) Delete(container *domain.Container) {
	force := true
	volumes := true

	options := &containers.RemoveOptions{
		Force:   &force,
		Volumes: &volumes,
	}

	_, err := containers.Remove(r.ctx, container.Id, options)
	if err != nil {
		panic(err)
	}
}

func (r *Runtime) Start(container *domain.Container) {
	err := containers.Start(r.ctx, container.Id, nil)
	if err != nil {
		panic(err)
	}
}

func (r *Runtime) Stop(container *domain.Container) {
	// err := containers.Stop(r.context, container.Id, nil)
	// if err != nil {
	// 	panic(err)
	// }

	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()
	ready := make(chan bool)
	done := make(chan bool)

	go func() {
		err := containers.Attach(r.ctx, container.Id, stdinReader, stdoutWriter, stdoutWriter, ready, nil)
		if err != nil {
			panic(err)
		}
		done <- true
	}()

	<-ready

	_, err := stdinWriter.Write([]byte("stop\n"))
	if err != nil {
		panic(err)
	}
	stdinWriter.Close()

	<-done

	stdinReader.Close()
	stdoutReader.Close()
	stdoutWriter.Close()
}

func (r *Runtime) Attach(container *domain.Container, stdin io.Reader, stdout io.Writer, stderr io.Writer, done chan struct{}) error {
	ready := make(chan bool)
	logs := make(chan string)

	go func() {
		defer close(logs)
		logTail := "30"
		options := &containers.LogOptions{
			Tail: &logTail,
		}
		err := containers.Logs(r.ctx, container.Id, options, logs, nil)
		if err != nil {
			panic(err)
		}
	}()

	for log := range logs {
		_, err := io.WriteString(stdout, log)
		if err != nil {
			panic(err)
		}
	}

	// TODO does this leak?
	go func() {
		defer close(done)

		err := containers.Attach(r.ctx, container.Id, stdin, stdout, stderr, ready, nil)
		if err != nil {
			panic(err)
		}
	}()
	<-ready
	return nil
}

func (r *Runtime) Events() chan domain.RuntimeEvent {
	events := make(chan domain.RuntimeEvent)

	stream := true
	options := &system.EventsOptions{
		Stream: &stream,
	}

	runtimeEvents := make(chan entities.Event)
	// runtimeEventsCancel := make(chan bool)

	err := system.Events(r.ctx, runtimeEvents, nil, options)
	if err != nil {
		panic(err)
	}

	go func() {
		for runtimeEvent := range runtimeEvents {
			containerId := runtimeEvent.Actor.ID
			serverId := runtimeEvent.Actor.Attributes["com.github.ayeama.panel.api.server.id"]
			if serverId == "" {
				continue
			}

			switch runtimeEvent.Type {
			case dockerevents.ContainerEventType:
				switch runtimeEvent.Action {
				case dockerevents.ActionCreate,
					dockerevents.ActionStart,
					dockerevents.ActionStop,
					"died":

					events <- domain.RuntimeEventServerStatusChanged{
						ServerId:    serverId,
						ContainerId: containerId,
					}
				}
			default:
			}
		}
	}()

	return events
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
