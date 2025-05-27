package runtime

import (
	"context"
	"fmt"
	"io"
	"math/rand"
	"sync"

	"github.com/ayeama/panel/api/internal/domain"
	nettypes "github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/system"
	"github.com/containers/podman/v5/pkg/domain/entities"
	"github.com/containers/podman/v5/pkg/specgen"
	dockerevents "github.com/docker/docker/api/types/events"
)

type PodmanNode struct {
	nodeName string
	nodeUri  string
	ctx      context.Context
}

type PodmanRuntime struct {
	mutex       sync.RWMutex
	connections map[string]*PodmanNode
	events      chan domain.Event
}

func NewPodmanRuntime() *PodmanRuntime {
	return &PodmanRuntime{
		connections: make(map[string]*PodmanNode, 0),
		events:      make(chan domain.Event),
	}
}

func (r *PodmanRuntime) AddNode(node *domain.Node) error {
	ctx, err := bindings.NewConnectionWithIdentity(context.Background(), node.Uri, "/home/alex/.ssh/id_rsa", false)
	if err != nil {
		panic(err)
	}

	r.mutex.Lock()
	defer r.mutex.Unlock()

	r.connections[node.Name] = &PodmanNode{nodeName: node.Name, nodeUri: node.Uri, ctx: ctx}

	return nil
}

func (r *PodmanRuntime) RemoveNode(name string) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	// TODO double check
	_, ok := r.connections[name]
	if !ok {
		return fmt.Errorf("failed to find node: %s", name)
	}

	delete(r.connections, name)

	return nil
}

func (r *PodmanRuntime) Node(name string) (Node, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	node, ok := r.connections[name]
	if !ok {
		return nil, fmt.Errorf("failed to find node: %s", name)
	}

	return node, nil
}

func (r *PodmanRuntime) RandomNode() (Node, error) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	keys := make([]string, 0, len(r.connections))
	for k := range r.connections {
		keys = append(keys, k)
	}
	name := keys[rand.Intn(len(keys))]

	node, ok := r.connections[name]
	if !ok {
		return nil, fmt.Errorf("failed to find node: %s", name)
	}

	return node, nil
}

func (r *PodmanRuntime) Events() <-chan domain.Event {
	return nil
}

func (r *PodmanRuntime) Ping() {
	var wg sync.WaitGroup
	for _, node := range r.connections {
		wg.Add(1)
		go func() {
			_, err := system.Info(node.ctx, nil)
			if err != nil {
				panic(err)
			}
			fmt.Printf("pinged node %s", node.Name)
		}()
	}
	wg.Wait()
}

func (r *PodmanRuntime) Close() {}

func (n *PodmanNode) Name() string {
	return n.nodeName
}

func (n *PodmanNode) Uri() string {
	return n.nodeUri
}

func (r *PodmanNode) Create(server *domain.Server) string {
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

	spec := specgen.NewSpecGenerator("localhost/server/minecraft:latest", false)
	// spec := specgen.NewSpecGenerator("localhost/server/valheim:latest", false)
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

	spec.Labels = make(map[string]string)
	spec.Labels["com.github.ayeama.panel.api.server.id"] = server.Id

	// todo container restart: unless-stopped

	resp, err := containers.CreateWithSpec(r.ctx, spec, nil)
	if err != nil {
		panic(err)
	}

	return resp.ID
}

func (r *PodmanNode) Delete(container *domain.Container) {
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

func (r *PodmanNode) Start(container *domain.Container) {
	err := containers.Start(r.ctx, container.Id, nil)
	if err != nil {
		panic(err)
	}
}

func (r *PodmanNode) Stop(container *domain.Container) {
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

func (r *PodmanNode) Stats(container *domain.Container) chan domain.ContainerStat {
	all := false
	stream := true
	interval := 1

	options := &containers.StatsOptions{
		All:      &all,
		Stream:   &stream,
		Interval: &interval,
	}

	resp, err := containers.Stats(r.ctx, []string{container.Id}, options)
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

func (r *PodmanNode) Attach(container *domain.Container, stdin io.Reader, stdout io.Writer, stderr io.Writer, done chan struct{}) error {
	resp, err := containers.Inspect(r.ctx, container.Id, nil)
	if err != nil {
		panic(err)
	}

	if resp.State.Status != "running" {
		return fmt.Errorf("container not running")
	}

	ready := make(chan bool)
	logs := make(chan string)

	go func() {
		defer close(logs)
		logTail := "30"
		options := &containers.LogOptions{
			Tail: &logTail,
		}
		err = containers.Logs(r.ctx, container.Id, options, logs, nil)
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

func (r *PodmanNode) Events() chan domain.Event {
	runtimeEvents := make(chan entities.Event)

	stream := true

	options := &system.EventsOptions{
		Stream: &stream,
	}

	err := system.Events(r.ctx, runtimeEvents, nil, options)
	if err != nil {
		panic(err)
	}

	events := make(chan domain.Event)

	go func() {
		for event := range runtimeEvents {
			serverId := event.Actor.ID
			containerId := event.Actor.Attributes["com.github.ayeama.panel.api.server.id"]

			if containerId == "" {
				continue
			}

			switch event.Type {
			case dockerevents.ContainerEventType:
				switch event.Action {
				case dockerevents.ActionCreate:
					events <- domain.ServerCreatedEvent{
						Id:          containerId,
						ContainerId: serverId,
					}
				case dockerevents.ActionStart:
					events <- domain.ServerStartedEvent{
						Id:          containerId,
						ContainerId: serverId,
					}
				case dockerevents.ActionStop, "died":
					events <- domain.ServerStoppedEvent{
						Id:          containerId,
						ContainerId: serverId,
					}
				case dockerevents.ActionRemove:
					events <- domain.ServerDeletedEvent{
						Id:          containerId,
						ContainerId: serverId,
					}
				default:
				}
			default:
			}
		}
	}()

	return events
}

func (r *PodmanNode) Status(container *domain.Container) string {
	resp, err := containers.Inspect(r.ctx, container.Id, nil)
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
