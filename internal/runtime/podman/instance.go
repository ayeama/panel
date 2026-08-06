package podman

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"path"
	osRuntime "runtime"
	"strings"

	"github.com/ayeama/panel/internal/runtime"
	"github.com/ayeama/panel/internal/types"
	"github.com/google/uuid"
	"github.com/opencontainers/runtime-spec/specs-go"
	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/bindings/volumes"
	entitiesTypes "go.podman.io/podman/v6/pkg/domain/entities/types"
	"go.podman.io/podman/v6/pkg/specgen"
)

// TODO pass resources limits into containers as environment variables for scripts?
func (r *Runtime) InstanceCreate(options types.InstanceCreate) (types.Instance, error) {
	// rimageID, err := r.imageID(imageID)
	// if err != nil {
	// 	return types.Instance{}, &runtime.Error{Op: "read", Resource: "image", ID: "", Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	// }

	image, err := images.GetImage(*r.ctx, options.Image, nil)
	if err != nil {
		return types.Instance{}, &runtime.Error{Op: "read", Resource: "image", ID: "", Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	spec := specgen.NewSpecGenerator(image.ID, false)

	spec.Volumes = make([]*specgen.NamedVolume, 0)
	for k := range image.Config.Volumes {
		volumeLabels := make(map[string]string)
		volumeLabels[instanceVolumeLabelID] = uuid.NewString()

		volumeCreateOptions := entitiesTypes.VolumeCreateOptions{
			Labels: volumeLabels,
		}

		volume, err := volumes.Create(*r.ctx, volumeCreateOptions, &volumes.CreateOptions{})
		if err != nil {
			return types.Instance{}, &runtime.Error{Op: "create", Resource: "instance", ID: "", Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
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

	if options.Resources.CPU > 0 {
		cpuPeriod := uint64(container_cpu_us)
		cpuQuota := int64(float64(cpuPeriod) * options.Resources.CPU)
		spec.ResourceLimits.CPU = &specs.LinuxCPU{
			Period: &cpuPeriod,
			Quota:  &cpuQuota,
		}
	}

	if options.Resources.Memory > 0 {
		memLimit := int64(options.Resources.Memory * container_memory_gb)
		spec.ResourceLimits.Memory = &specs.LinuxMemory{
			Limit: &memLimit,
		}
	}

	if options.Resources.Disk > 0 {
		// TODO not implemented
	}

	spec.Labels = make(map[string]string)

	id := uuid.NewString()
	spec.Labels[instanceLabelID] = id
	spec.Labels[instanceLabelWebhooks] = strings.Join(options.Webhooks, ",")

	container, err := containers.CreateWithSpec(*r.ctx, spec, nil)
	if err != nil {
		return types.Instance{}, &runtime.Error{Op: "create", Resource: "instance", ID: "", Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	err = containers.ContainerInit(*r.ctx, container.ID, nil)
	if err != nil {
		return types.Instance{}, &runtime.Error{Op: "create", Resource: "instance", ID: "", Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	instance, err := r.InstanceRead(id)
	if err != nil {
		return types.Instance{}, err
	}

	return instance, nil
}

func (r *Runtime) InstanceRead(id string) (types.Instance, error) {
	filters := map[string][]string{"label": {instanceLabelID + "=" + id}}
	containerListOptions := containers.ListOptions{}
	containerListOptions.WithAll(true).WithFilters(filters)

	containerList, err := containers.List(*r.ctx, &containerListOptions)
	if err != nil {
		return types.Instance{}, &runtime.Error{Op: "read", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	for _, container := range containerList {
		instanceID := container.Labels[instanceLabelID]
		if instanceID == id {
			containerDeep, err := containers.Inspect(*r.ctx, container.ID, nil)
			if err != nil {
				return types.Instance{}, &runtime.Error{Op: "read", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
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
			if container.Labels[instanceLabelWebhooks] != "" {
				webhooks = strings.Split(container.Labels[instanceLabelWebhooks], ",")
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
					CPU:    cpus,
					Memory: memory,
					Disk:   disk,
				},
				Webhooks: webhooks,
			}
			return instance, nil
		}
	}

	return types.Instance{}, &runtime.Error{Op: "read", Resource: "instance", ID: id, Err: runtime.ErrNotFound}
}

func (r *Runtime) InstanceReadMany() ([]types.Instance, error) {
	filters := map[string][]string{"label": {instanceLabelID}}
	containerListOptions := containers.ListOptions{}
	containerListOptions.WithAll(true).WithFilters(filters)

	containerList, err := containers.List(*r.ctx, &containerListOptions)
	if err != nil {
		return []types.Instance{}, &runtime.Error{Op: "read", Resource: "instance", ID: "", Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	instances := make([]types.Instance, 0)
	for _, container := range containerList {
		instanceID := container.Labels[instanceLabelID]

		if instanceID == "" {
			continue
		}

		containerDeep, err := containers.Inspect(*r.ctx, container.ID, nil)
		if err != nil {
			return []types.Instance{}, &runtime.Error{Op: "read", Resource: "instance", ID: "", Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
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

		// TODO handle empty webhook string ""?
		webhooks := strings.Split(container.Labels[instanceLabelWebhooks], ",")

		instances = append(instances, types.Instance{
			ID:     instanceID,
			Name:   container.Names[0],
			Image:  container.Image,
			Status: container.State,
			Ports:  transformPorts(container.Ports),
			Resources: types.InstanceResources{
				CPU:    cpus,
				Memory: memory,
				Disk:   disk,
			},
			Webhooks: webhooks,
		})
	}

	return instances, nil
}

func (r *Runtime) InstanceDelete(id string) error {
	containerID, err := r.containerID(id)
	if err != nil {
		return &runtime.Error{Op: "read", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	containerRemoveOptions := containers.RemoveOptions{}
	containerRemoveOptions.WithForce(true).WithVolumes(true).WithTimeout(1)

	_, err = containers.Remove(*r.ctx, containerID, &containerRemoveOptions)
	if err != nil {
		return &runtime.Error{Op: "delete", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	return nil
}

func (r *Runtime) InstanceStart(id string) error {
	containerID, err := r.containerID(id)
	if err != nil {
		return &runtime.Error{Op: "read", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	if err = containers.Start(*r.ctx, containerID, nil); err != nil {
		return &runtime.Error{Op: "start", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	return nil
}

func (r *Runtime) InstanceStop(id string) error {
	containerID, err := r.containerID(id)
	if err != nil {
		return &runtime.Error{Op: "stop", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	timeout := uint(1)
	containerStopOptions := containers.StopOptions{
		Timeout: &timeout,
	}

	if err = containers.Stop(*r.ctx, containerID, &containerStopOptions); err != nil {
		return &runtime.Error{Op: "stop", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	return nil
}

func (r *Runtime) InstanceAttach(id string, stdin io.Reader, stdout io.Writer, stderr io.Writer, ready chan bool) error {
	containerID, err := r.containerID(id)
	if err != nil {
		return &runtime.Error{Op: "attach", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	containerAttachOptions := containers.AttachOptions{}
	containerAttachOptions.WithDetachKeys("")
	if err = containers.Attach(*r.ctx, containerID, stdin, stdout, stderr, ready, &containerAttachOptions); err != nil {
		if errors.Is(err, io.ErrClosedPipe) {
			return nil
		}
		return &runtime.Error{Op: "attach", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	return nil
}

// TODO not working
func (r *Runtime) InstanceStats(id string, stats chan types.InstanceStat) error {
	containerID, err := r.containerID(id)
	if err != nil {
		return &runtime.Error{Op: "stats", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	containerDeep, err := containers.Inspect(*r.ctx, containerID, nil)
	if err != nil {
		return &runtime.Error{Op: "stats", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	cpus := 0.0

	if containerDeep.HostConfig.CpuPeriod > 0 && containerDeep.HostConfig.CpuQuota > 0 {
		cpus = float64(containerDeep.HostConfig.CpuQuota) / float64(containerDeep.HostConfig.CpuPeriod)
	}

	if cpus == 0 {
		cpus = float64(osRuntime.NumCPU())
	}

	containerStatsOptions := containers.StatsOptions{}
	containerStatsOptions.WithInterval(1)

	statsReport, err := containers.Stats(*r.ctx, []string{containerDeep.ID}, &containerStatsOptions)
	if err != nil {
		return &runtime.Error{Op: "stats", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	for statReports := range statsReport {
		for _, stat := range statReports.Stats {
			cpu := stat.CPU / cpus

			netTX := uint64(0)
			netRX := uint64(0)
			for _, net := range stat.Network {
				netTX += net.TxBytes
				netRX += net.RxBytes
			}

			stats <- types.InstanceStat{
				CPUPercent:     cpu,
				MemoryPercent:  stat.MemPerc,
				DiskPercent:    float64(0), // TODO disk usage
				NetworkTXBytes: netTX,
				NetworkRXBytes: netRX,
			}
		}
	}

	return nil
}

func (r *Runtime) InstanceBackup(id string, manifest *types.InstanceBackupManifest, zw *zip.Writer) error {
	containerID, err := r.containerID(id)
	if err != nil {
		return &runtime.Error{Op: "backup", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	containerDeep, err := containers.Inspect(*r.ctx, containerID, nil)
	if err != nil {
		return &runtime.Error{Op: "backup", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	for _, mount := range containerDeep.Mounts {
		volume, err := volumes.Inspect(*r.ctx, mount.Name, nil)
		if err != nil {
			return &runtime.Error{Op: "backup", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
		}

		id := volume.Labels[instanceVolumeLabelID]

		w, err := zw.Create(fmt.Sprintf("volumes/%s", id))
		if err != nil {
			return &runtime.Error{Op: "backup", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
		}

		err = volumes.Export(*r.ctx, volume.Name, w)
		if err != nil {
			return &runtime.Error{Op: "backup", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
		}

		manifest.Mounts = append(manifest.Mounts, types.InstanceBackupManifestMount{
			ID:          id,
			Destination: mount.Destination,
		})
	}

	return nil
}

func (r *Runtime) InstanceRestore(id string, manifest *types.InstanceBackupManifest, zr *zip.Reader) error {
	// TODO don't rely on manifest? user could malform it as an attack
	for _, manifestMount := range manifest.Mounts {
		volumeName, err := r.volumeName(manifestMount.ID)
		if err != nil {
			return &runtime.Error{Op: "restore", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
		}

		for _, filePath := range zr.File {
			if path.Base(filePath.Name) == manifestMount.ID {
				fmt.Println("found a match:", volumeName, filePath.Name, manifestMount.ID)
				f, err := zr.Open(filePath.Name)
				if err != nil {
					return &runtime.Error{Op: "restore", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
				}
				defer f.Close()

				err = volumes.Import(*r.ctx, volumeName, f)
				if err != nil {
					return &runtime.Error{Op: "restore", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
				}
				break
			}
			fmt.Println("did not find a match:", volumeName, filePath.Name, manifestMount.ID)
		}

	}

	return nil
}

func (r *Runtime) InstanceLogs(id string, logs chan string) error {
	containerID, err := r.containerID(id)
	if err != nil {
		return &runtime.Error{Op: "logs", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	containerLogOptions := containers.LogOptions{}
	containerLogOptions.WithStderr(true).WithStdout(true).WithTimestamps(true)
	if err = containers.Logs(*r.ctx, containerID, &containerLogOptions, logs, logs); err != nil {
		return &runtime.Error{Op: "logs", Resource: "instance", ID: id, Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	return nil
}
