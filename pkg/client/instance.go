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
	ID        string
	Name      string
	Image     string
	Status    string
	Ports     map[string]string
	Resources InstanceResources

	Webhooks []string
}

type InstanceCreate struct {
	Image     string
	Resources InstanceResources
	Webhooks  []string
}

type InstanceResources struct {
	CPU    float64
	Memory float64
	Disk   float64
}

type InstanceStat struct {
	CPUPercent     float64
	MemoryPercent  float64
	DiskPercent    float64
	NetworkTxBytes uint64
	NetworkRxBytes uint64
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

func (s *InstanceService) Delete(ctx context.Context, id string) error {
	url := fmt.Sprintf("/instances/%s", id)
	req, err := s.client.newRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return err
	}

	_, err = s.client.http.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (s *InstanceService) Start(ctx context.Context, id string) error {
	url := fmt.Sprintf("/instances/%s/start", id)
	req, err := s.client.newRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	_, err = s.client.http.Do(req)
	if err != nil {
		return err
	}

	return nil
}

func (s *InstanceService) Stop(ctx context.Context, id string) error {
	url := fmt.Sprintf("/instances/%s/stop", id)
	req, err := s.client.newRequestWithContext(ctx, http.MethodPost, url, nil)
	if err != nil {
		return err
	}

	_, err = s.client.http.Do(req)
	if err != nil {
		return err
	}

	return nil
}
