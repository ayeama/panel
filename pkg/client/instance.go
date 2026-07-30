package panel

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
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

type InstanceResources struct {
	Cpu    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
	Disk   float64 `json:"disk"`
}

type InstanceStat struct {
	CpuPercent     float64 `json:"cpu_percent"`
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

func (s *InstanceService) Create(ctx context.Context) error {
	// TODO
	return nil
}

func (s *InstanceService) Read(ctx context.Context) ([]Instance, error) {
	instances := make([]Instance, 0)

	req, err := s.client.newRequestWithContext(ctx, http.MethodGet, "/instances", nil)
	if err != nil {
		return instances, err
	}

	resp, err := s.client.http.Do(req)
	if err != nil {
		return instances, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&instances); err != nil {
		return instances, err
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

	if err = json.NewDecoder(resp.Body).Decode(&instance); err != nil {
		return instance, err
	}

	return instance, nil
}
