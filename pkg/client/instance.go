package panel

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ayeama/panel/pkg/api"
)

type Instance struct {
	ID        string            `json:"id"`
	Name      string            `json:"name"`
	Image     string            `json:"image"`
	Status    string            `json:"status"`
	Ports     map[string]string `json:"ports"`
	Resources InstanceResources `json:"resources"`

	Webhooks []string `json:"webhooks"`
}

type InstanceCreate struct {
	Image     string
	Resources InstanceResources
	Webhooks  []string
}

type InstanceResources struct {
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
	Disk   float64 `json:"disk"`
}

type InstanceStat struct {
	CPUPercent     float64 `json:"cpu_percent"`
	MemoryPercent  float64 `json:"memory_percent"`
	DiskPercent    float64 `json:"disk_percent"`
	NetworkTxBytes uint64  `json:"network_tx_bytes"`
	NetworkRxBytes uint64  `json:"network_rx_bytes"`
}

type InstanceBackupManifest struct {
	Instance // TODO don't need all members

	Version string                        `json:"version"`
	Mounts  []InstanceBackupManifestMount `json:"volumes"`
}

type InstanceBackupManifestMount struct {
	ID          string `json:"id"`
	Destination string `json:"destination"`
}

type InstanceService service

func (s *InstanceService) Create(ctx context.Context, options InstanceCreate) (Instance, error) {
	var instance Instance

	body, err := json.Marshal(api.InstanceCreateRequest{
		Image: options.Image,
		Resources: api.InstanceResources{
			CPU:    options.Resources.CPU,
			Memory: options.Resources.Memory,
			Disk:   options.Resources.Disk,
		},
		Webhooks: options.Webhooks,
	})
	if err != nil {
		return instance, err
	}

	req, err := s.client.newRequestWithContext(ctx, http.MethodPost, "/instances", bytes.NewReader(body))
	if err != nil {
		return instance, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := s.client.http.Do(req)
	if err != nil {
		return instance, err
	}
	defer resp.Body.Close()

	var respInstance api.Instance
	if err = json.NewDecoder(resp.Body).Decode(&respInstance); err != nil {
		return instance, err
	}

	instance = Instance{
		ID:     respInstance.ID,
		Name:   respInstance.Name,
		Image:  respInstance.Image,
		Status: respInstance.Status,
		Ports:  respInstance.Ports,
		Resources: InstanceResources{
			CPU:    respInstance.Resources.CPU,
			Memory: respInstance.Resources.Memory,
			Disk:   respInstance.Resources.Disk,
		},
		Webhooks: respInstance.Webhooks,
	}

	return instance, nil
}

func (s *InstanceService) Read(ctx context.Context) ([]Instance, error) {
	var instances []Instance

	req, err := s.client.newRequestWithContext(ctx, http.MethodGet, "/instances", nil)
	if err != nil {
		return instances, err
	}

	resp, err := s.client.http.Do(req)
	if err != nil {
		return instances, err
	}
	defer resp.Body.Close()

	var respInstances []api.Instance
	if err = json.NewDecoder(resp.Body).Decode(&respInstances); err != nil {
		return instances, err
	}

	instances = make([]Instance, len(respInstances))
	for i, respInstance := range respInstances {
		instances[i] = Instance{
			ID:     respInstance.ID,
			Name:   respInstance.Name,
			Image:  respInstance.Image,
			Status: respInstance.Status,
			Ports:  respInstance.Ports,
			Resources: InstanceResources{
				CPU:    respInstance.Resources.CPU,
				Memory: respInstance.Resources.Memory,
				Disk:   respInstance.Resources.Disk,
			},
			Webhooks: respInstance.Webhooks,
		}
	}

	return instances, nil
}

func (s *InstanceService) ReadOne(ctx context.Context, id string) (Instance, error) {
	var instance Instance

	url := fmt.Sprintf("/instances/%s", id)
	req, err := s.client.newRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return instance, err
	}

	resp, err := s.client.http.Do(req)
	if err != nil {
		return instance, err
	}
	defer resp.Body.Close()

	var respInstance api.Instance
	if err = json.NewDecoder(resp.Body).Decode(&respInstance); err != nil {
		return instance, err
	}

	instance = Instance{
		ID:     respInstance.ID,
		Name:   respInstance.Name,
		Image:  respInstance.Image,
		Status: respInstance.Status,
		Ports:  respInstance.Ports,
		Resources: InstanceResources{
			CPU:    respInstance.Resources.CPU,
			Memory: respInstance.Resources.Memory,
			Disk:   respInstance.Resources.Disk,
		},
		Webhooks: respInstance.Webhooks,
	}

	return instance, nil
}
