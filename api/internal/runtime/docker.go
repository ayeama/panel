package runtime

import (
	"context"
	"io"
	"strconv"
	"strings"
	"sync"

	"github.com/ayeama/panel/api/internal/domain"
	dockercontainer "github.com/docker/docker/api/types/container"
	dockerevents "github.com/docker/docker/api/types/events"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
)

type Docker struct {
	client *client.Client
}

func (r Docker) New() (Runtime, error) {
	client, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return &Docker{}, err
	}
	return &Docker{client: client}, nil
}

func (r *Docker) inspect(id string) dockercontainer.InspectResponse {
	resp, err := r.client.ContainerInspect(context.Background(), id)
	if err != nil {
		panic(err)
	}
	return resp
}

func (r *Docker) Name(id string) string {
	name := r.inspect(id).Name
	name = strings.TrimPrefix(name, "/")
	return name
}

func (r *Docker) Status(id string) string {
	return r.inspect(id).State.Status
}

func (r *Docker) Port(id string) string {
	resp := r.inspect(id)
	for _, bindings := range resp.HostConfig.PortBindings {
		for _, binding := range bindings {
			return binding.HostPort
		}
	}
	return ""
}

func (r *Docker) Running(id string) bool {
	return r.inspect(id).State.Running
}

func (r *Docker) Create(id string, image string) string {
	exposedPorts := nat.PortSet{
		"25565/tcp": struct{}{},
	}

	config := &dockercontainer.Config{
		Image:        image,
		AttachStdin:  true,
		OpenStdin:    true,
		Tty:          true,
		Labels:       map[string]string{"com.github.ayeama.panel.api.server.id": id},
		ExposedPorts: exposedPorts,
	}

	port, err := freeHostPort()
	if err != nil {
		panic(err)
	}
	hostPort := strconv.FormatUint(uint64(port), 10)

	hostConfig := &dockercontainer.HostConfig{
		PortBindings: nat.PortMap{
			"25565/tcp": []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: hostPort},
			},
		},
	}

	var networkConfig *network.NetworkingConfig

	resp, err := r.client.ContainerCreate(context.Background(), config, hostConfig, networkConfig, nil, "")
	if err != nil {
		panic(err)
	}

	return resp.ID
}

func (r *Docker) Delete(container *domain.Container) {
	options := dockercontainer.RemoveOptions{
		RemoveVolumes: true,
		Force:         true,
	}

	err := r.client.ContainerRemove(context.Background(), container.Id, options)
	if err != nil {
		// do nothing? log warning?
	}
}

func (r *Docker) Start(container *domain.Container) {
	err := r.client.ContainerStart(context.Background(), container.Id, dockercontainer.StartOptions{})
	if err != nil {
		panic(err)
	}
}

func (r *Docker) Stop(container *domain.Container) {
	// TODO attach and send the shutdown command
	options := dockercontainer.StopOptions{}
	err := r.client.ContainerStop(context.Background(), container.Id, options)
	if err != nil {
		panic(err)
	}
}

func (r *Docker) Attach(container *domain.Container, stdin io.Reader, stdout io.Writer, stderr io.Writer, done chan struct{}) error {
	ctx := context.Background()
	var once sync.Once

	options := dockercontainer.AttachOptions{
		Stream: true,
		Stdin:  true,
		Stdout: true,
		Stderr: true,
	}
	resp, err := r.client.ContainerAttach(ctx, container.Id, options)
	if err != nil {
		panic(err)
	}

	go func() {
		<-done
		resp.Close()
	}()

	go func() {
		_, err := io.Copy(resp.Conn, stdin)
		if err != nil {
			panic(err)
		}
		once.Do(func() { close(done) })
	}()

	go func() {
		_, err := io.Copy(stdout, resp.Reader)
		if err != nil {
			panic(err)
		}
		once.Do(func() { close(done) })
	}()

	return nil
}

func (r *Docker) Events() chan domain.RuntimeEvent {
	events := make(chan domain.RuntimeEvent)

	runtimeEvents, _ := r.client.Events(context.Background(), dockerevents.ListOptions{})
	// if err != nil {
	// 	panic(err)
	// }

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
