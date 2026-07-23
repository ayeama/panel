package handler

import (
	"archive/zip"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"time"

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
	mux.HandleFunc("GET /instances/{id}/backup", h.handleInstanceBackup)
	mux.HandleFunc("POST /instances/{id}/restore", h.handleInstanceRestore)
	mux.HandleFunc("GET /instances/{id}/logs", h.handleInstanceLogs)
}

func (h *InstanceHandler) handleInstanceCreate(w http.ResponseWriter, r *http.Request) {
	type instanceCreateRequest struct {
		ImageID   string                  `json:"image_id"`
		Resources types.InstanceResources `json:"resources"`
		Webhooks  []string                `json:"webhooks"`
	}
	var req instanceCreateRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Fatal(err)
	}

	instance, err := h.runtime.InstanceCreate(req.ImageID, req.Resources, req.Webhooks)
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

	w.WriteHeader(http.StatusNoContent)
}

func (h *InstanceHandler) handleInstanceStop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	if err := h.runtime.InstanceStop(id); err != nil {
		log.Fatal(err)
	}

	w.WriteHeader(http.StatusNoContent)
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
		if err := h.runtime.InstanceAttach(id, stdinReader, stdoutWriter, stderrWriter, ready); err != nil {
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

			select {
			case msgs <- msg:
			case <-ctx.Done():
				return
			}
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

			select {
			case msgs <- msg:
			case <-ctx.Done():
				return
			}
		}
	}()

	select {
	case <-ready:
	case <-ctx.Done():
		return
	}

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

	// TODO add context with cancel?

	id := r.PathValue("id")

	stats := make(chan types.InstanceStat)

	go func() {
		if err := h.runtime.InstanceStats(id, stats); err != nil {
			log.Println("about to fail in stats instanestats")
			log.Fatal(err)
		}
	}()

	for stat := range stats {
		if err := c.WriteJSON(stat); err != nil {
			return
		}
	}
}

func (h *InstanceHandler) handleInstanceBackup(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	instance, err := h.runtime.InstanceRead(id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(
			"attachment; filename=\"panel-%s-%s-backup.zip\"",
			instance.Name,
			time.Now().Format("20060102150405"),
		),
	)
	w.WriteHeader(http.StatusOK)

	zw := zip.NewWriter(w)
	defer zw.Close()

	manifest := types.InstanceBackupManifest{
		Instance: instance,
		Version:  "1",
	}

	if err = h.runtime.InstanceBackup(instance.ID, &manifest, zw); err != nil {
		log.Fatal(err)
	}

	f, err := zw.Create("manifest.json")
	if err != nil {
		log.Fatal(err)
	}

	if err = json.NewEncoder(f).Encode(manifest); err != nil {
		log.Fatal(err)
	}
}

func (h *InstanceHandler) handleInstanceRestore(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	instance, err := h.runtime.InstanceRead(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	file, _, err := r.FormFile("backup")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	tmpFile, err := os.CreateTemp("", "panel-backup-*.zip")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile.Name())
	defer tmpFile.Close()

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	stat, err := tmpFile.Stat()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	zr, err := zip.NewReader(tmpFile, stat.Size())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO trim the zip filename prefix showhow?
	var manifestFile fs.File
	for _, f := range zr.File {
		if path.Base(f.Name) == "manifest.json" {
			manifestFile, err = zr.Open(f.Name)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			break
		}
	}
	if manifestFile == nil {
		http.Error(w, "manifest.json not found", http.StatusBadRequest)
		return
	}
	defer manifestFile.Close()

	var manifest types.InstanceBackupManifest
	err = json.NewDecoder(manifestFile).Decode(&manifest)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if manifest.ID != instance.ID {
		http.Error(w, "backup manifest ID does not match instance ID", http.StatusBadRequest)
		return
	}

	err = h.runtime.InstanceRestore(instance.ID, &manifest, zr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError) // TODO
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *InstanceHandler) handleInstanceLogs(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	id := r.PathValue("id")

	instance, err := h.runtime.InstanceRead(id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set(
		"Content-Disposition",
		fmt.Sprintf(
			"attachment; filename=\"panel-%s-%s-logs.zip\"",
			instance.Name,
			time.Now().Format("20060102150405"),
		),
	)
	w.WriteHeader(http.StatusOK)

	logs := make(chan string)

	go func() {
		defer close(logs)

		if err := h.runtime.InstanceLogs(instance.ID, logs); err != nil {
			cancel()
		}
	}()

	zw := zip.NewWriter(w)
	defer zw.Close()

	f, err := zw.Create("logs.txt")
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case log, ok := <-logs:
			if !ok {
				return
			}
			if _, err := fmt.Fprint(f, log); err != nil {
				return
			}
		}
	}
}
