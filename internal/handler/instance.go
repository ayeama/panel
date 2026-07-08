package handler

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type InstanceHandler struct {
	podman *context.Context
}

func NewInstanceHandler(podman *context.Context) InstanceHandler {
	return InstanceHandler{podman}
}

func (h *InstanceHandler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("POST /instances", handleInstanceCreate)
	mux.HandleFunc("GET /instances", handleInstanceReadMany)
	mux.HandleFunc("GET /instances/{id}", handleInstanceRead)
	mux.HandleFunc("PUT /instances/{id}", handleInstanceUpdate)
	mux.HandleFunc("DELETE /instances/{id}", handleInstanceDelete)
	mux.HandleFunc("GET /instances/{id}/attach", handleInstanceAttach)
	mux.HandleFunc("GET /instances/{id}/stats", handleInstanceStats)
}

func handleInstanceCreate(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
}

func handleInstanceReadMany(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "[{\"id\":\"1\",\"name\":\"another\",\"image\":\"minecraft\",\"status\":\"stopped\"},{\"id\":\"2\",\"name\":\"one\",\"image\":\"valheim\",\"status\":\"running\"}]")
	w.Header().Add("Content-Type", "application/json")
}

func handleInstanceRead(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "{\"id\":\"1\",\"name\":\"another\",\"image\":\"minecraft\",\"status\":\"stopped\"}")
	w.Header().Add("Content-Type", "application/json")
}

func handleInstanceUpdate(w http.ResponseWriter, r *http.Request) {}

func handleInstanceDelete(w http.ResponseWriter, r *http.Request) {}

func handleInstanceAttach(w http.ResponseWriter, r *http.Request) {
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

func handleInstanceStats(w http.ResponseWriter, r *http.Request) {
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
