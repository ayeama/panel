package domain

import (
	"context"
	"io"

	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/specgen"
)

type Pod struct {
	Id string
}

type Server struct {
	Id     string
	Name   string
	Status string
	Pod    Pod
}

func (s *Server) Create() {
	conn, err := bindings.NewConnection(context.Background(), "unix:///run/user/1000/podman/podman.sock")
	if err != nil {
		panic(err)
	}

	stdin := true
	terminal := true

	spec := specgen.NewSpecGenerator("localhost/server/minecraft:latest", false)
	spec.Stdin = &stdin
	spec.Terminal = &terminal

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
