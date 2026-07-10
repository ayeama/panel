package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/ayeama/panel/internal/types"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.podman.io/podman/v6/pkg/bindings/containers"
	"go.podman.io/podman/v6/pkg/bindings/images"
	"go.podman.io/podman/v6/pkg/specgen"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type InstanceHandler struct {
	runtime *context.Context
}

func NewInstanceHandler(runtime *context.Context) InstanceHandler {
	return InstanceHandler{runtime}
}

func (h *InstanceHandler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("POST /instances", h.handleInstanceCreate)
	mux.HandleFunc("GET /instances", h.handleInstanceReadMany)
	mux.HandleFunc("GET /instances/{id}", h.handleInstanceRead)
	mux.HandleFunc("PUT /instances/{id}", h.handleInstanceUpdate)
	mux.HandleFunc("DELETE /instances/{id}", h.handleInstanceDelete)
	mux.HandleFunc("POST /instances/{id}/start", h.handleInstanceStart)
	mux.HandleFunc("POST /instances/{id}/stop", h.handleInstanceStop)
	mux.HandleFunc("GET /instances/{id}/attach", h.handleInstanceAttach)
	mux.HandleFunc("GET /instances/{id}/stats", h.handleInstanceStats)
}

func (h *InstanceHandler) handleInstanceCreate(w http.ResponseWriter, r *http.Request) {
	type instanceCreateRequest struct {
		ImageID string `json:"image_id"`
	}
	var request instanceCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Fatal(err)
	}

	filters := map[string][]string{"label": {"com.github.ayeama.panel.image.id=" + request.ImageID}}
	imageListOptions := &images.ListOptions{}
	imageListOptions.WithAll(false).WithFilters(filters)

	imageList, err := images.List(*h.runtime, imageListOptions)
	if err != nil {
		log.Fatal(err)
	}

	imageID := ""
	for _, image := range imageList {
		id := image.Labels["com.github.ayeama.panel.image.id"]
		if id == request.ImageID {
			imageID = image.ID
			break
		}
	}

	// TODO create volume ourselves

	spec := specgen.NewSpecGenerator(imageID, false)

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

	uuid, err := uuid.NewUUID()
	if err != nil {
		log.Fatal(err)
	}
	spec.Labels["com.github.ayeama.panel.instance.id"] = uuid.String()

	container, err := containers.CreateWithSpec(*h.runtime, spec, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = containers.ContainerInit(*h.runtime, container.ID, nil)
	if err != nil {
		log.Fatal(err)
	}

	type instanceCreateResponse struct {
		InstanceID string `json:"instance_id"`
	}

	containerResponse := instanceCreateResponse{
		InstanceID: uuid.String(),
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(containerResponse); err != nil {
		log.Fatal(err)
	}
}

func (h *InstanceHandler) handleInstanceReadMany(w http.ResponseWriter, r *http.Request) {
	filters := map[string][]string{"label": {"com.github.ayeama.panel.instance.id"}}
	containerListOptions := containers.ListOptions{}
	containerListOptions.WithAll(true).WithFilters(filters)

	containerList, err := containers.List(*h.runtime, &containerListOptions)
	if err != nil {
		log.Fatal(err)
	}

	var instanceResponse []types.Instance
	for _, instance := range containerList {
		id := instance.Labels["com.github.ayeama.panel.instance.id"]

		if id == "" {
			continue
		}

		instanceResponse = append(instanceResponse, types.Instance{
			ID:     id,
			Name:   instance.Names[0],
			Image:  instance.Image,
			Status: instance.State,
		})
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(instanceResponse); err != nil {
		log.Fatal(err)
	}
}

func (h *InstanceHandler) handleInstanceRead(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("id")

	filters := map[string][]string{"label": {"com.github.ayeama.panel.instance.id=" + instanceID}}
	containerListOptions := containers.ListOptions{}
	containerListOptions.WithAll(true).WithFilters(filters)

	containerList, err := containers.List(*h.runtime, &containerListOptions)
	if err != nil {
		log.Fatal(err)
	}

	var instanceResponse types.Instance
	for _, instance := range containerList {
		id := instance.Labels["com.github.ayeama.panel.instance.id"]
		if id == instanceID {
			instanceResponse = types.Instance{
				ID:     id,
				Name:   instance.Names[0],
				Image:  instance.Image,
				Status: instance.State,
			}
		}
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(instanceResponse); err != nil {
		log.Fatal(err)
	}
}

func (h *InstanceHandler) handleInstanceUpdate(w http.ResponseWriter, r *http.Request) {}

func (h *InstanceHandler) handleInstanceDelete(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("id")

	filters := map[string][]string{"label": {"com.github.ayeama.panel.instance.id=" + instanceID}}
	containerListOptions := containers.ListOptions{}
	containerListOptions.WithAll(true).WithFilters(filters)

	containerList, err := containers.List(*h.runtime, &containerListOptions)
	if err != nil {
		log.Fatal(err)
	}

	containerRemoveOptions := containers.RemoveOptions{}
	containerRemoveOptions.WithForce(true).WithVolumes(true).WithTimeout(1)

	for _, instance := range containerList {
		id := instance.Labels["com.github.ayeama.panel.instance.id"]
		if id == instanceID {
			containers.Remove(*h.runtime, instance.ID, &containerRemoveOptions)
			break
		}
	}
}

func (h *InstanceHandler) handleInstanceStart(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("id")

	all := true
	filters := map[string][]string{"label": {"com.github.ayeama.panel.instance.id=" + instanceID}}
	containerListOptions := containers.ListOptions{
		All:     &all,
		Filters: filters,
	}

	containerList, err := containers.List(*h.runtime, &containerListOptions)
	if err != nil {
		log.Fatal(err)
	}

	for _, instance := range containerList {
		id := instance.Labels["com.github.ayeama.panel.instance.id"]
		if id == instanceID {
			err := containers.Start(*h.runtime, instance.ID, nil)
			if err != nil {
				log.Fatal(err)
			}
			break
		}
	}
}

func (h *InstanceHandler) handleInstanceStop(w http.ResponseWriter, r *http.Request) {
	instanceID := r.PathValue("id")

	all := true
	filters := map[string][]string{"label": {"com.github.ayeama.panel.instance.id=" + instanceID}}
	containerListOptions := containers.ListOptions{
		All:     &all,
		Filters: filters,
	}

	containerList, err := containers.List(*h.runtime, &containerListOptions)
	if err != nil {
		log.Fatal(err)
	}

	timeout := uint(1)
	containerStopOptions := containers.StopOptions{
		Timeout: &timeout,
	}

	for _, instance := range containerList {
		id := instance.Labels["com.github.ayeama.panel.instance.id"]
		if id == instanceID {
			err := containers.Stop(*h.runtime, instance.ID, &containerStopOptions)
			if err != nil {
				log.Fatal(err)
			}
			break
		}
	}
}

func (h *InstanceHandler) handleInstanceAttach(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, msg, err := c.ReadMessage()
		if err != nil {
			return
		}

		if err := c.WriteMessage(mt, msg); err != nil {
			return
		}
	}
}

func (h *InstanceHandler) handleInstanceStats(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		cpu := rand.Float64() * 100
		memory := rand.Float64() * 100
		disk := rand.Float64() * 100

		msg := []byte(fmt.Sprintf("{\"cpu\":%.2f,\"memory\":%.2f,\"disk\":%.2f}", cpu, memory, disk))
		if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
			return
		}
		time.Sleep(time.Second)
	}
}
