package podman

import (
	"context"
	"errors"
	"strconv"

	"go.podman.io/podman/v6/pkg/bindings"
	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/bindings/volumes"

	netTypes "go.podman.io/common/libnetwork/types"
)

const (
	container_cpu_us    = 1_000_000
	container_memory_gb = 1_000_000_000

	imageLabelID          string = "com.github.ayeama.panel.image.id"
	imageLabelVersion     string = "com.github.ayeama.panel.image.version"
	imageLabelName        string = "com.github.ayeama.panel.image.name"
	imageLabelDescription string = "com.github.ayeama.panel.image.description"

	instanceLabelID       string = "com.github.ayeama.panel.instance.id"
	instanceLabelWebhooks string = "com.github.ayeama.panel.instance.webhooks"

	instanceVolumeLabelID string = "com.github.ayeama.panel.volume.id"
)

type Runtime struct {
	ctx *context.Context
}

type Config struct {
	URI string
}

// TODO: pass in context?
func New(config *Config) (*Runtime, error) {
	c, err := bindings.NewConnection(context.Background(), config.URI)
	if err != nil {
		return nil, err
	}
	return &Runtime{&c}, nil
}

func (r *Runtime) imageID(id string) (string, error) {
	filters := map[string][]string{"label": {imageLabelID + "=" + id}}
	imageListOptions := &images.ListOptions{}
	imageListOptions.WithAll(false).WithFilters(filters)

	imageList, err := images.List(*r.ctx, imageListOptions)
	if err != nil {
		return "", err
	}

	for _, image := range imageList {
		imageID := image.Labels[imageLabelID]
		if imageID == id {
			return image.ID, nil
		}
	}

	return "", errors.New("image not found")
}

func (r *Runtime) containerID(id string) (string, error) {
	filters := map[string][]string{"label": {instanceLabelID + "=" + id}}
	containerListOptions := containers.ListOptions{}
	containerListOptions.WithAll(true).WithFilters(filters)

	containerList, err := containers.List(*r.ctx, &containerListOptions)
	if err != nil {
		return "", err
	}

	for _, container := range containerList {
		instanceID := container.Labels[instanceLabelID]
		if instanceID == id {
			return container.ID, nil
		}
	}

	return "", errors.New("container not found")
}

func (r *Runtime) volumeName(id string) (string, error) {
	filters := map[string][]string{"label": {instanceVolumeLabelID + "=" + id}}
	volumeListOptions := volumes.ListOptions{}
	volumeListOptions.WithFilters(filters)

	volumeList, err := volumes.List(*r.ctx, &volumeListOptions)
	if err != nil {
		return "", err
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
