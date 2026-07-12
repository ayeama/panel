package runtime

import (
	"context"
	"errors"
	"io"
	"log"
	"strconv"

	"github.com/ayeama/panel/internal/types"
	"github.com/google/uuid"
	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/bindings/system"
	"go.podman.io/podman/v6/pkg/bindings/volumes"
	entitiesTypes "go.podman.io/podman/v6/pkg/domain/entities/types"
	"go.podman.io/podman/v6/pkg/specgen"

	mobyEvents "github.com/moby/moby/api/types/events"

	netTypes "go.podman.io/common/libnetwork/types"
)

type PodmanRuntime struct {
	ctx *context.Context
}

func NewPodmanRuntime(ctx *context.Context) PodmanRuntime {
	return PodmanRuntime{ctx}
}

func (r *PodmanRuntime) ImageRead(id string) (types.Image, error) {
	filters := map[string][]string{"label": {types.ImageLabelID + "=" + id}}
	imageListOptions := &images.ListOptions{}
	imageListOptions.WithAll(false).WithFilters(filters)

	imageList, err := images.List(*r.ctx, imageListOptions)
	if err != nil {
		log.Fatal(err)
	}

	for _, image := range imageList {
		imageID := image.Labels[types.ImageLabelID]
		if imageID == id {
			return types.Image{
				ID:   imageID,
				Name: image.Names[0],
			}, nil
		}
	}

	return types.Image{}, errors.New("image not found")
}

func (r *PodmanRuntime) ImageReadMany() ([]types.Image, error) {
	filters := map[string][]string{"label": {types.ImageLabelID}}
	imageListOptions := &images.ListOptions{}
	imageListOptions.WithAll(false).WithFilters(filters)

	imageList, err := images.List(*r.ctx, imageListOptions)
	if err != nil {
		log.Fatal(err)
	}

	images := make([]types.Image, 0)
	for _, image := range imageList {
		id := image.Labels[types.ImageLabelID]

		if id == "" {
			continue
		}

		if len(image.Names) < 1 {
			continue
		}

		images = append(images, types.Image{
			ID:   image.Labels[types.ImageLabelID],
			Name: image.Names[0],
		})
	}
	return images, nil
}

func (r *PodmanRuntime) InstanceCreate(imageID string) (types.Instance, error) {
	// TODO replace imageName with actual image id not panel id

	rimageID, err := r.imageID(imageID)
	if err != nil {
		log.Fatal(err)
	}

	image, err := images.GetImage(*r.ctx, rimageID, nil)
	if err != nil {
		log.Fatal(err)
	}

	spec := specgen.NewSpecGenerator(image.ID, false)

	spec.Volumes = make([]*specgen.NamedVolume, 0, len(image.Config.Volumes))
	for k := range image.Config.Volumes {
		volumeLabels := make(map[string]string)
		volumeLabels["com.github.ayeama.panel.volume.id"] = uuid.NewString()

		volumeCreateOptions := entitiesTypes.VolumeCreateOptions{
			Labels: volumeLabels,
		}

		volume, err := volumes.Create(*r.ctx, volumeCreateOptions, &volumes.CreateOptions{})
		if err != nil {
			log.Fatal(err)
		}

		spec.Volumes = append(spec.Volumes, &specgen.NamedVolume{
			Name: volume.Name,
			Dest: k,
		})
	}

	publish := true
	spec.PublishExposedPorts = &publish

	stdin := true
	spec.Stdin = &stdin

	terminal := true
	spec.Terminal = &terminal

	// cpus := 1.0
	// mem := 1.0
	// cpuPeriod := uint64(100000)
	// cpuQuota := int64(float64(cpuPeriod) * cpus)
	// memLimit := int64(mem * 1000000000)
	// spec.ResourceLimits = &specs.LinuxResources{
	// 	CPU: &specs.LinuxCPU{
	// 		Period: &cpuPeriod,
	// 		Quota:  &cpuQuota,
	// 	},
	// 	Memory: &specs.LinuxMemory{
	// 		Limit: &memLimit,
	// 	},
	// }

	spec.Labels = make(map[string]string)

	id := uuid.NewString()
	spec.Labels[types.InstanceLabelID] = id
	spec.Labels[types.InstanceLabelWebhook] = "http://localhost:8001/webhook" // TODO

	container, err := containers.CreateWithSpec(*r.ctx, spec, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = containers.ContainerInit(*r.ctx, container.ID, nil)
	if err != nil {
		log.Fatal(err)
	}

	instance, err := r.InstanceRead(id)
	if err != nil {
		log.Fatal(err)
	}

	return instance, nil
}

func (r *PodmanRuntime) InstanceRead(id string) (types.Instance, error) {
	filters := map[string][]string{"label": {types.InstanceLabelID + "=" + id}}
	containerListOptions := containers.ListOptions{}
	containerListOptions.WithAll(true).WithFilters(filters)

	containerList, err := containers.List(*r.ctx, &containerListOptions)
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range containerList {
		instanceID := container.Labels[types.InstanceLabelID]
		if instanceID == id {
			instance := types.Instance{
				ID:      instanceID,
				Name:    container.Names[0],
				Image:   container.Image,
				Status:  container.State,
				Ports:   transformPorts(container.Ports),
				Webhook: container.Labels[types.InstanceLabelWebhook],
			}
			return instance, nil
		}
	}

	return types.Instance{}, errors.New("instance not found")
}

func (r *PodmanRuntime) InstanceReadMany() ([]types.Instance, error) {
	filters := map[string][]string{"label": {types.InstanceLabelID}}
	containerListOptions := containers.ListOptions{}
	containerListOptions.WithAll(true).WithFilters(filters)

	containerList, err := containers.List(*r.ctx, &containerListOptions)
	if err != nil {
		log.Fatal(err)
	}

	var instances []types.Instance
	for _, container := range containerList {
		instanceID := container.Labels[types.InstanceLabelID]

		if instanceID == "" {
			continue
		}

		instances = append(instances, types.Instance{
			ID:      instanceID,
			Name:    container.Names[0],
			Image:   container.Image,
			Status:  container.State,
			Ports:   transformPorts(container.Ports),
			Webhook: container.Labels[types.InstanceLabelWebhook],
		})
	}

	return instances, nil
}

func (r *PodmanRuntime) InstanceDelete(id string) error {
	containerID, err := r.containerID(id)
	if err != nil {
		log.Fatal(err)
	}

	containerRemoveOptions := containers.RemoveOptions{}
	containerRemoveOptions.WithForce(true).WithVolumes(true).WithTimeout(1)

	_, err = containers.Remove(*r.ctx, containerID, &containerRemoveOptions)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func (r *PodmanRuntime) InstanceStart(id string) error {
	containerID, err := r.containerID(id)
	if err != nil {
		log.Fatal(err)
	}

	if err = containers.Start(*r.ctx, containerID, nil); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (r *PodmanRuntime) InstanceStop(id string) error {
	containerID, err := r.containerID(id)
	if err != nil {
		log.Fatal(err)
	}

	timeout := uint(1)
	containerStopOptions := containers.StopOptions{
		Timeout: &timeout,
	}

	if err = containers.Stop(*r.ctx, containerID, &containerStopOptions); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (r *PodmanRuntime) InstanceAttach(id string, stdin io.Reader, stdout io.Writer, stderr io.Writer, ready chan bool) error {
	containerID, err := r.containerID(id)
	if err != nil {
		log.Fatal(err)
	}

	// msgs := make(chan string, 55)

	// go func() {
	// 	defer close(msgs)

	// 	containerLogOptions := containers.LogOptions{}
	// 	containerLogOptions.WithStdout(true).WithStderr(true).WithTail("50").WithFollow(false)
	// 	if err := containers.Logs(*r.ctx, containerID, &containerLogOptions, msgs, msgs); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }()

	// for msg := range msgs {
	// 	if _, err = stdout.Write([]byte(msg)); err != nil {
	// 		log.Fatal(err)
	// 	}
	// }

	containerAttachOptions := containers.AttachOptions{}
	containerAttachOptions.WithDetachKeys("")
	if err = containers.Attach(*r.ctx, containerID, stdin, stdout, stderr, ready, &containerAttachOptions); err != nil {
		log.Fatal(err)
	}

	return nil
}

// TODO not working
func (r *PodmanRuntime) InstanceStats(id string, stats chan types.InstanceStat) error {
	containerID, err := r.containerID(id)
	if err != nil {
		log.Fatal(err)
	}

	containerStatsOptions := containers.StatsOptions{}
	containerStatsOptions.WithInterval(1)

	statsReport, err := containers.Stats(*r.ctx, []string{containerID}, &containerStatsOptions)
	if err != nil {
		log.Fatal(err)
	}

	for statReports := range statsReport {
		for _, stat := range statReports.Stats {
			netTx := uint64(0)
			netRx := uint64(0)
			for _, net := range stat.Network {
				netTx += net.TxBytes
				netRx += net.RxBytes
			}

			stats <- types.InstanceStat{
				CpuPercent:     stat.CPU,
				MemoryPercent:  stat.MemPerc,
				DiskPercent:    float64(0), // TODO disk usage
				NetworkTxBytes: netTx,
				NetworkRxBytes: netRx,
			}
		}
	}

	return nil
}

func (r *PodmanRuntime) Events(events chan types.Event, cancel chan bool) error {
	podmanEvents := make(chan entitiesTypes.Event)
	if err := system.Events(*r.ctx, podmanEvents, cancel, nil); err != nil {
		log.Fatal(err)
	}

	for podmanEvent := range podmanEvents {
		switch podmanEvent.Type {
		case mobyEvents.ContainerEventType:
			switch podmanEvent.Action {
			case mobyEvents.ActionCreate:
				events <- types.Event{
					Type:   types.EventTypeInstance,
					Action: types.EventActionCreate,
					Actor: types.Actor{
						ID:         podmanEvent.Actor.ID,
						Attributes: podmanEvent.Actor.Attributes,
					},
				}
			case mobyEvents.ActionRemove:
				events <- types.Event{
					Type:   types.EventTypeInstance,
					Action: types.EventActionDelete,
					Actor: types.Actor{
						ID:         podmanEvent.Actor.ID,
						Attributes: podmanEvent.Actor.Attributes,
					},
				}
			}
		default:
			break
		}
	}

	return nil
}

func (r *PodmanRuntime) imageID(id string) (string, error) {
	filters := map[string][]string{"label": {types.ImageLabelID + "=" + id}}
	imageListOptions := &images.ListOptions{}
	imageListOptions.WithAll(false).WithFilters(filters)

	imageList, err := images.List(*r.ctx, imageListOptions)
	if err != nil {
		log.Fatal(err)
	}

	for _, image := range imageList {
		imageID := image.Labels[types.ImageLabelID]
		if imageID == id {
			return image.ID, nil
		}
	}

	return "", errors.New("image not found")
}

func (r *PodmanRuntime) containerID(id string) (string, error) {
	filters := map[string][]string{"label": {types.InstanceLabelID + "=" + id}}
	containerListOptions := containers.ListOptions{}
	containerListOptions.WithAll(true).WithFilters(filters)

	containerList, err := containers.List(*r.ctx, &containerListOptions)
	if err != nil {
		log.Fatal(err)
	}

	for _, container := range containerList {
		instanceID := container.Labels[types.InstanceLabelID]
		if instanceID == id {
			return container.ID, nil
		}
	}

	return "", errors.New("instance not found")
}

func transformPorts(ports []netTypes.PortMapping) map[string]string {
	transPorts := map[string]string{}
	for _, port := range ports {
		containerPort := strconv.FormatUint(uint64(port.ContainerPort), 10)
		hostPort := strconv.FormatUint(uint64(port.HostPort), 10)

		if transPorts[containerPort] != "" {
			continue
		}

		transPorts[containerPort] = hostPort
	}

	return transPorts
}
