package runtime

import (
	"context"
	"fmt"
	"io"

	"github.com/ayeama/panel/api/internal/config"
	"github.com/ayeama/panel/api/internal/domain"
	nettypes "github.com/containers/common/libnetwork/types"
	"github.com/containers/podman/v5/pkg/api/handlers"
	"github.com/containers/podman/v5/pkg/bindings"
	"github.com/containers/podman/v5/pkg/bindings/containers"
	"github.com/containers/podman/v5/pkg/bindings/system"
	"github.com/containers/podman/v5/pkg/bindings/volumes"
	"github.com/containers/podman/v5/pkg/domain/entities"
	entitiesTypes "github.com/containers/podman/v5/pkg/domain/entities/types"
	"github.com/containers/podman/v5/pkg/specgen"
	dockerContainer "github.com/docker/docker/api/types/container"
	dockerevents "github.com/docker/docker/api/types/events"
)

type Podman struct {
	ctx context.Context
}

func (r Podman) New() (Runtime, error) {
	ctx, err := bindings.NewConnection(context.Background(), config.Config.RuntimeUri)
	if err != nil {
		return &Podman{}, err
	}
	return &Podman{ctx: ctx}, nil
}

func (r *Podman) Inspect(container_id string) domain.Container {
	resp, err := containers.Inspect(r.ctx, container_id, nil)
	if err != nil {
		panic(err)
	}
	container := domain.Container{
		Id:     resp.ID,
		Name:   resp.Name,
		Status: resp.State.Status,
		Ports:  make([]string, 0),
	}
	for _, ports := range resp.NetworkSettings.Ports {
		// TODO im assuming there is only ever one port, it is a list though
		container.Ports = append(container.Ports, fmt.Sprintf("%s:%s", config.Config.ServerHost, ports[0].HostPort))
	}
	return container
}

func (r *Podman) Create(id string, tag string) string {
	volumeOptions := entitiesTypes.VolumeCreateOptions{}
	volumeResponse, err := volumes.Create(r.ctx, volumeOptions, nil)
	if err != nil {
		panic(err)
	}

	volumes := make([]*specgen.NamedVolume, 0)
	volume := &specgen.NamedVolume{
		Name: volumeResponse.Name,
		Dest: "/data",
	}
	volumes = append(volumes, volume)

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

	spec := specgen.NewSpecGenerator(tag, false)
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
	spec.Volumes = volumes
	spec.PortMappings = portMappings

	// restartRetries := uint(3)
	// spec.RestartPolicy = "unless-stopped"
	// spec.RestartRetries = &restartRetries

	spec.Labels = make(map[string]string)
	spec.Labels["com.github.ayeama.panel.api.server.id"] = id

	serverResponse, err := containers.CreateWithSpec(r.ctx, spec, nil)
	if err != nil {
		panic(err)
	}

	err = containers.ContainerInit(r.ctx, serverResponse.ID, nil)
	if err != nil {
		panic(err)
	}

	return serverResponse.ID
}

func (r *Podman) Delete(container_id string) {
	force := true
	volumes := true
	timeout := uint(1)

	options := &containers.RemoveOptions{
		Force:   &force,
		Volumes: &volumes,
		Timeout: &timeout,
	}

	_, err := containers.Remove(r.ctx, container_id, options)
	if err != nil {
		panic(err)
	}
}

func (r *Podman) Start(container_id string) {
	err := containers.Start(r.ctx, container_id, nil)
	if err != nil {
		panic(err)
	}
}

func (r *Podman) Stop(container_id string) {
	timeout := uint(1)
	options := &containers.StopOptions{Timeout: &timeout}
	err := containers.Stop(r.ctx, container_id, options)
	if err != nil {
		panic(err)
	}

	// resp := r.inspect(container_id)
	// sftpId := resp.Config.Labels["com.github.ayeama.panel.api.server.sftp.id"]
	// sftpStopTimeout := uint(1)
	// sftpStopOptions := &containers.StopOptions{Timeout: &sftpStopTimeout}
	// err := containers.Stop(r.ctx, sftpId, sftpStopOptions)
	// if err != nil {
	// 	panic(err)
	// }

	// stdinReader, stdinWriter := io.Pipe()
	// stdoutReader, stdoutWriter := io.Pipe()
	// ready := make(chan bool)
	// done := make(chan bool)

	// go func() {
	// 	err := containers.Attach(r.ctx, container_id, stdinReader, stdoutWriter, stdoutWriter, ready, nil)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	done <- true
	// }()

	// <-ready

	// writeDone := make(chan error, 1)
	// timeout := time.Second * 30

	// go func() {
	// 	_, err := stdinWriter.Write([]byte("stop\n"))
	// 	writeDone <- err
	// }()

	// select {
	// case err := <-writeDone:
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	stdinWriter.Close()
	// case <-time.After(timeout):
	// 	_ = stdinWriter.CloseWithError(fmt.Errorf("write timeout after %s", timeout))
	// 	err := containers.Stop(r.ctx, container_id, nil)
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// }

	// <-done

	// _ = stdinReader.Close()
	// _ = stdoutReader.Close()
	// _ = stdoutWriter.Close()
}

func (r *Podman) Attach(container_id string, stdin io.Reader, stdout io.Writer, stderr io.Writer, done chan struct{}) error {
	ready := make(chan bool)
	logs := make(chan string)

	go func() {
		defer close(logs)
		logTail := "30"
		options := &containers.LogOptions{
			Tail: &logTail,
		}
		err := containers.Logs(r.ctx, container_id, options, logs, nil)
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

		err := containers.Attach(r.ctx, container_id, stdin, stdout, stderr, ready, nil)
		if err != nil {
			panic(err)
		}
	}()
	<-ready
	return nil
}

func (r *Podman) Events() chan domain.RuntimeEvent {
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

func (r *Podman) CreateSidecar(id string, tag string, server_id string) string {
	server, err := containers.Inspect(r.ctx, server_id, nil)
	if err != nil {
		panic(err)
	}

	// get volumes mounted to the server
	volumes := make([]*specgen.NamedVolume, 0)
	for _, mount := range server.Mounts {
		volumes = append(volumes, &specgen.NamedVolume{
			Name: mount.Name,
			Dest: fmt.Sprintf("/data%s", mount.Destination),
		})
	}

	port, err := freeHostPort()
	if err != nil {
		panic(err)
	}

	spec := specgen.NewSpecGenerator(tag, false)
	spec.Volumes = volumes
	spec.PortMappings = []nettypes.PortMapping{{HostPort: port, ContainerPort: 22}}
	spec.RestartPolicy = "always"

	resp, err := containers.CreateWithSpec(r.ctx, spec, nil)
	if err != nil {
		panic(err)
	}

	return resp.ID
}

func (r *Podman) InjectCredentials(container_id string) {
	credentials := []string{
		"",
		"",
	}

	for _, credential := range credentials {
		createOptions := &handlers.ExecCreateConfig{
			ExecOptions: dockerContainer.ExecOptions{
				Cmd: []string{"add-key", credential},
			},
		}
		exec, err := containers.ExecCreate(r.ctx, container_id, createOptions)
		if err != nil {
			panic(err)
		}

		err = containers.ExecStart(r.ctx, exec, nil)
		if err != nil {
			panic(err)
		}
	}

}
