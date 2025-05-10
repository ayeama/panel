package domain

import (
	"context"
	"fmt"
	"io"

	nettypes "github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
	"github.com/opencontainers/runtime-spec/specs-go"
)

type Pod struct {
	Id     string
	Status string
}

type Server struct {
	Id   string
	Name string
	Pod  Pod
}

func (s *Server) Create() {
	conn, err := bindings.NewConnection(context.Background(), "unix:///run/user/1000/podman/podman.sock")
	if err != nil {
		panic(err)
	}

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

	createResponse, err := containers.CreateWithSpec(conn, spec, nil)
	if err != nil {
		panic(err)
	}

	s.Pod.Id = createResponse.ID

	s.Start()
}

func (s *Server) Start() {
	conn, err := bindings.NewConnection(context.Background(), "unix:///run/user/1000/podman/podman.sock")
	if err != nil {
		panic(err)
	}

	err = containers.Start(conn, s.Pod.Id, nil)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Stop() {
	conn, err := bindings.NewConnection(context.Background(), "unix:///run/user/1000/podman/podman.sock")
	if err != nil {
		panic(err)
	}

	err = containers.Stop(conn, s.Pod.Id, nil)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Remove() {
	conn, err := bindings.NewConnection(context.Background(), "unix:///run/user/1000/podman/podman.sock")
	if err != nil {
		panic(err)
	}

	force := true
	volumes := true

	options := &containers.RemoveOptions{
		Force:   &force,
		Volumes: &volumes,
	}

	removeResponse, err := containers.Remove(conn, s.Pod.Id, options)
	if err != nil {
		panic(err)
	}

	fmt.Println(removeResponse)
}

func (s *Server) Logs(stdout chan string, stderr chan string) {
	conn, err := bindings.NewConnection(context.Background(), "unix:///run/user/1000/podman/podman.sock")
	if err != nil {
		panic(err)
	}

	tail := "25"
	options := &containers.LogOptions{
		Tail: &tail,
	}

	err = containers.Logs(conn, s.Pod.Id, options, stdout, stderr)
	if err != nil {
		panic(err)
	}
}

func (s *Server) Attach(stdin io.Reader, stdout io.Writer, stderr io.Writer) {
	conn, err := bindings.NewConnection(context.Background(), "unix:///run/user/1000/podman/podman.sock")
	if err != nil {
		panic(err)
	}

	attachReady := make(chan bool)

	go func() {
		err = containers.Attach(conn, s.Pod.Id, stdin, stdout, stderr, attachReady, nil)
		if err != nil {
			panic(err)
		}
	}()

	<-attachReady
}

func (s *Server) Inspect() {
	conn, err := bindings.NewConnection(context.Background(), "unix:///run/user/1000/podman/podman.sock")
	if err != nil {
		panic(err)
	}

	inspectResponse, err := containers.Inspect(conn, s.Pod.Id, nil)
	if err != nil {
		panic(err)
	}

	s.Pod.Status = inspectResponse.State.Status
}
