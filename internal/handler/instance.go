package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ayeama/panel/internal/runtime"
	"github.com/ayeama/panel/internal/types"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type InstanceHandler struct {
	runtime runtime.Runtime
}

func NewInstanceHandler(runtime runtime.Runtime) InstanceHandler {
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

	// TODO
	// mux.HandleFunc("GET /instances/{id}/backup", h.handleInstanceBackup)
	// mux.HandleFunc("POST /instances/{id}/restore", h.handleInstanceRestore)
	// mux.HandleFunc("GET /instances/{id}/logs", h.handleInstanceLogs)
}

func (h *InstanceHandler) handleInstanceCreate(w http.ResponseWriter, r *http.Request) {
	type instanceCreateRequest struct {
		ImageID string `json:"image_id"`
	}
	var request instanceCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		log.Fatal(err)
	}

	image, err := h.runtime.ImageRead(request.ImageID)
	if err != nil {
		log.Fatal(err)
	}

	instance, err := h.runtime.InstanceCreate(image.Name)
	if err != nil {
		log.Fatal(err)
	}

	type instanceCreateResponse struct {
		InstanceID string `json:"instance_id"`
	}

	containerResponse := instanceCreateResponse{
		InstanceID: instance.ID,
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(containerResponse); err != nil {
		log.Fatal(err)
	}
}

func (h *InstanceHandler) handleInstanceReadMany(w http.ResponseWriter, r *http.Request) {
	instances, err := h.runtime.InstanceReadMany()
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(instances); err != nil {
		log.Fatal(err)
	}
}

func (h *InstanceHandler) handleInstanceRead(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	instance, err := h.runtime.InstanceRead(id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(instance); err != nil {
		log.Fatal(err)
	}
}

func (h *InstanceHandler) handleInstanceUpdate(w http.ResponseWriter, r *http.Request) {}

func (h *InstanceHandler) handleInstanceDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.runtime.InstanceDelete(id); err != nil {
		log.Fatal(err)
	}
}

func (h *InstanceHandler) handleInstanceStart(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.runtime.InstanceStart(id); err != nil {
		log.Fatal(err)
	}
}

func (h *InstanceHandler) handleInstanceStop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.runtime.InstanceStop(id); err != nil {
		log.Fatal(err)
	}
}

func (h *InstanceHandler) handleInstanceAttach(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()

	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := r.PathValue("id")

	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()
	stderrReader, stderrWriter := io.Pipe()

	ready := make(chan bool)

	go func() {
		defer stdinReader.Close()
		defer stdoutWriter.Close()
		defer stderrWriter.Close()

		if err = h.runtime.InstanceAttach(id, stdinReader, stdoutWriter, stderrWriter, ready); err != nil {
			cancel()
		}
	}()

	msgs := make(chan []byte)

	go func() {
		defer cancel()
		defer stdinWriter.Close()

		for {
			_, msg, err := c.ReadMessage()
			if err != nil {
				return
			}

			_, err = stdinWriter.Write(msg)
			if err != nil {
				return
			}
		}
	}()

	go func() {
		defer cancel()
		defer stdoutReader.Close()

		buf := make([]byte, 1024)
		for {
			n, err := stdoutReader.Read(buf)
			if err != nil {
				return
			}

			msg := make([]byte, n)
			copy(msg, buf[:n])
			msgs <- msg
		}
	}()

	go func() {
		defer cancel()
		defer stderrReader.Close()

		buf := make([]byte, 1024)
		for {
			n, err := stderrReader.Read(buf)
			if err != nil {
				return
			}

			msg := make([]byte, n)
			copy(msg, buf[:n])
			msgs <- msg
		}
	}()

	<-ready

	for {
		select {
		case <-ctx.Done():
			return
		case msg := <-msgs:
			if err := c.WriteMessage(websocket.TextMessage, msg); err != nil {
				return
			}
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

	id := r.PathValue("id")

	stats := make(chan types.InstanceStat)

	// TODO handle lifecycle
	go func() {
		if err = h.runtime.InstanceStats(id, stats); err != nil {
			log.Fatal(err)
		}
	}()

	for stat := range stats {
		if err = c.WriteJSON(stat); err != nil {
			log.Fatal(err)
		}
	}
}
