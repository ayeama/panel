package runtime

import (
	"archive/zip"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"path"
	osruntime "runtime"
	"strconv"
	"strings"

	"github.com/ayeama/panel/internal/types"
	"github.com/google/uuid"
	"github.com/opencontainers/runtime-spec/specs-go"
	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/bindings/system"
	"go.podman.io/podman/v6/pkg/bindings/volumes"
	entitiesTypes "go.podman.io/podman/v6/pkg/domain/entities/types"
	"go.podman.io/podman/v6/pkg/specgen"

	mobyEvents "github.com/moby/moby/api/types/events"

	netTypes "go.podman.io/common/libnetwork/types"
)

const (
	container_cpu_us    = 1_00_000
	container_memory_gb = 1_000_000_000

	// TODO: better place?
	instanceVolumeLabelID string = "com.github.ayeama.panel.volume.id"
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

// TODO pass resources limits into containers as environment variables for scripts?
func (r *PodmanRuntime) InstanceCreate(imageID string, resources types.InstanceResources, webhooks []string) (types.Instance, error) {
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
		volumeLabels[instanceVolumeLabelID] = uuid.NewString()

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

	spec.ResourceLimits = &specs.LinuxResources{}

	if resources.Cpu > 0 {
		cpuPeriod := uint64(container_cpu_us)
		cpuQuota := int64(float64(cpuPeriod) * resources.Cpu)
		spec.ResourceLimits.CPU = &specs.LinuxCPU{
			Period: &cpuPeriod,
			Quota:  &cpuQuota,
		}
	}

	if resources.Memory > 0 {
		memLimit := int64(resources.Memory * container_memory_gb)
		spec.ResourceLimits.Memory = &specs.LinuxMemory{
			Limit: &memLimit,
		}
	}

	if resources.Disk > 0 {
		// TODO not implemented
	}

	spec.Labels = make(map[string]string)

	id := uuid.NewString()
	spec.Labels[types.InstanceLabelID] = id
	spec.Labels[types.InstanceLabelWebhooks] = strings.Join(webhooks, ",")

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
			containerDeep, err := containers.Inspect(*r.ctx, container.ID, nil)
			if err != nil {
				log.Fatal(err)
			}

			cpus := 0.0
			memory := 0.0
			disk := 0.0 // TODO

			if containerDeep.HostConfig.CpuPeriod > 0 && containerDeep.HostConfig.CpuQuota > 0 {
				cpus = float64(containerDeep.HostConfig.CpuQuota) / float64(containerDeep.HostConfig.CpuPeriod)
			}

			if containerDeep.HostConfig.Memory > 0 {
				memory = float64(containerDeep.HostConfig.Memory) / container_memory_gb
			}

			var webhooks []string
			if container.Labels[types.InstanceLabelWebhooks] != "" {
				webhooks = strings.Split(container.Labels[types.InstanceLabelWebhooks], ",")
			} else {
				webhooks = make([]string, 0)
			}

			instance := types.Instance{
				ID:     instanceID,
				Name:   container.Names[0],
				Image:  container.Image,
				Status: container.State,
				Ports:  transformPorts(container.Ports),
				Resources: types.InstanceResources{
					Cpu:    cpus,
					Memory: memory,
					Disk:   disk,
				},
				Webhooks: webhooks,
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

	instances := make([]types.Instance, 0)
	for _, container := range containerList {
		instanceID := container.Labels[types.InstanceLabelID]

		if instanceID == "" {
			continue
		}

		containerDeep, err := containers.Inspect(*r.ctx, container.ID, nil)
		if err != nil {
			log.Fatal(err)
		}

		cpus := 0.0
		memory := 0.0
		disk := 0.0 // TODO

		if containerDeep.HostConfig.CpuPeriod > 0 && containerDeep.HostConfig.CpuQuota > 0 {
			cpus = float64(containerDeep.HostConfig.CpuQuota) / float64(containerDeep.HostConfig.CpuPeriod)
		}

		if containerDeep.HostConfig.Memory > 0 {
			memory = float64(containerDeep.HostConfig.Memory) / container_memory_gb
		}

		webhooks := strings.Split(container.Labels[types.InstanceLabelWebhooks], ",")

		instances = append(instances, types.Instance{
			ID:     instanceID,
			Name:   container.Names[0],
			Image:  container.Image,
			Status: container.State,
			Ports:  transformPorts(container.Ports),
			Resources: types.InstanceResources{
				Cpu:    cpus,
				Memory: memory,
				Disk:   disk,
			},
			Webhooks: webhooks,
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
	// 	containerLogOptions.WithStderr(true).WithStdout(true).WithTail("50").WithFollow(false)
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
		if errors.Is(err, io.ErrClosedPipe) {
			return nil
		}

		fmt.Println("about to error here:", err, err.Error())
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

	containerDeep, err := containers.Inspect(*r.ctx, containerID, nil)
	if err != nil {
		log.Fatal(err)
	}

	cpus := 0.0

	if containerDeep.HostConfig.CpuPeriod > 0 && containerDeep.HostConfig.CpuQuota > 0 {
		cpus = float64(containerDeep.HostConfig.CpuQuota) / float64(containerDeep.HostConfig.CpuPeriod)
	}

	if cpus == 0 {
		cpus = float64(osruntime.NumCPU())
	}

	containerStatsOptions := containers.StatsOptions{}
	containerStatsOptions.WithInterval(1)

	statsReport, err := containers.Stats(*r.ctx, []string{containerDeep.ID}, &containerStatsOptions)
	if err != nil {
		// TODO bug: 2026/07/15 09:14:57 write tcp [::1]:8000->[::1]:51528: write: broken pipe
		log.Println("about to fail in runtime stats stats")
		log.Fatal(err)
	}

	for statReports := range statsReport {
		for _, stat := range statReports.Stats {
			cpu := stat.CPU / cpus

			netTx := uint64(0)
			netRx := uint64(0)
			for _, net := range stat.Network {
				netTx += net.TxBytes
				netRx += net.RxBytes
			}

			stats <- types.InstanceStat{
				CpuPercent:     cpu,
				MemoryPercent:  stat.MemPerc,
				DiskPercent:    float64(0), // TODO disk usage
				NetworkTxBytes: netTx,
				NetworkRxBytes: netRx,
			}
		}
	}

	return nil
}

func (r *PodmanRuntime) InstanceBackup(id string, manifest *types.InstanceBackupManifest, zw *zip.Writer) error {
	containerID, err := r.containerID(id)
	if err != nil {
		log.Fatal(err)
	}

	containerDeep, err := containers.Inspect(*r.ctx, containerID, nil)
	if err != nil {
		log.Fatal(err)
	}

	for _, mount := range containerDeep.Mounts {
		volume, err := volumes.Inspect(*r.ctx, mount.Name, nil)
		if err != nil {
			log.Fatal(err)
		}

		id := volume.Labels[instanceVolumeLabelID]

		w, err := zw.Create(fmt.Sprintf("volumes/%s", id))
		if err != nil {
			log.Fatal(err)
		}

		err = volumes.Export(*r.ctx, volume.Name, w)
		if err != nil {
			log.Fatal(err)
		}

		manifest.Mounts = append(manifest.Mounts, types.InstanceBackupManifestMount{
			ID:          id,
			Destination: mount.Destination,
		})
	}

	return nil
}

func (r *PodmanRuntime) InstanceRestore(id string, manifest *types.InstanceBackupManifest, zr *zip.Reader) error {
	// containerID, err := r.containerID(id)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// containerDeep, err := containers.Inspect(*r.ctx, containerID, nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// TODO don't rely on manifest? user could malform it as an attack

	for _, manifestMount := range manifest.Mounts {
		volumeName, err := r.volumeName(manifestMount.ID)
		if err != nil {
			log.Println("about to fail runtime instance restore volume name")
			log.Fatal(err)
		}

		for _, filePath := range zr.File {
			if path.Base(filePath.Name) == manifestMount.ID {
				fmt.Println("found a match:", volumeName, filePath.Name, manifestMount.ID)
				f, err := zr.Open(filePath.Name)
				if err != nil {
					log.Fatal(err)
				}
				defer f.Close()

				err = volumes.Import(*r.ctx, volumeName, f)
				if err != nil {
					log.Fatal(err)
				}
				break
			}
			fmt.Println("did not find a match:", volumeName, filePath.Name, manifestMount.ID)
		}

	}

	return nil
}

func (r *PodmanRuntime) InstanceLogs(id string, logs chan string) error {
	containerID, err := r.containerID(id)
	if err != nil {
		log.Fatal(err)
	}

	containerLogOptions := containers.LogOptions{}
	containerLogOptions.WithStderr(true).WithStdout(true).WithTimestamps(true)
	if err = containers.Logs(*r.ctx, containerID, &containerLogOptions, logs, logs); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (r *PodmanRuntime) Events(events chan types.Event, cancel chan bool) error {
	podmanEvents := make(chan entitiesTypes.Event)
	if err := system.Events(*r.ctx, podmanEvents, cancel, nil); err != nil {
		log.Println("about to fail in runtime events")
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

func (r *PodmanRuntime) volumeName(id string) (string, error) {
	filters := map[string][]string{"label": {instanceVolumeLabelID + "=" + id}}
	volumeListOptions := volumes.ListOptions{}
	volumeListOptions.WithFilters(filters)

	volumeList, err := volumes.List(*r.ctx, &volumeListOptions)
	if err != nil {
		log.Fatal(err)
	}

	for _, volume := range volumeList {
		volumeName := volume.Labels[instanceVolumeLabelID]
		if volumeName == id {
			return volume.Name, nil
		}
	}

	return "", errors.New("volume not found")
}

func transformPorts(ports []netTypes.PortMapping) map[string]string {
	transPorts := make(map[string]string, 0)

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
