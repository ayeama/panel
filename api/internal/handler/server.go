package handler

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/ayeama/panel/api/internal/config"
	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/service"
	"github.com/ayeama/panel/api/pkg/api/types"
)

type ServerHandler struct {
	service *service.ServerService
}

func NewServerHandler(service *service.ServerService) *ServerHandler {
	return &ServerHandler{service: service}
}

func (h *ServerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var serverCreate types.ServerCreateRequest
	ReadRequestJson(r.Body, &serverCreate)

	// TODO request model validation
	if serverCreate.Image == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	server := h.service.Create(serverCreate.Image)

	output := types.ServerResponse{Id: server.Id, Name: server.Name, Status: server.Status, Address: fmt.Sprintf("%s:%s", config.Config.ServerHost, server.Container.Port)}
	WriteResponseJson(w, 200, output)
}

func (h *ServerHandler) Read(w http.ResponseWriter, r *http.Request) {
	pagination := domain.NewPagination(r.URL.Query())

	// TODO add helper function for conversion
	domainServerPaginated := h.service.Read(pagination)
	serverPaginated := types.PaginationResponse[types.ServerResponse]{
		Limit:  domainServerPaginated.Limit,
		Offset: domainServerPaginated.Offset,
		Total:  domainServerPaginated.Total,
		Items:  make([]types.ServerResponse, 0),
	}

	for _, s := range domainServerPaginated.Items {
		serverPaginated.Items = append(serverPaginated.Items, types.ServerResponse{Id: s.Id, Name: s.Name, Status: s.Status, Address: fmt.Sprintf("%s:%s", config.Config.ServerHost, s.Container.Port)})
	}

	WriteResponseJson(w, 200, serverPaginated)
}

func (h *ServerHandler) ReadOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	server := &domain.Server{Id: id}
	err := h.service.ReadOne(server)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			panic(err)
		}
	}

	output := types.ServerResponse{Id: server.Id, Name: server.Name, Status: server.Status, Address: fmt.Sprintf("%s:%s", config.Config.ServerHost, server.Container.Port)}
	WriteResponseJson(w, 200, output)
}

func (h *ServerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	server := &domain.Server{Id: id}
	h.service.Delete(server)

	w.WriteHeader(http.StatusNoContent)
}

func (h *ServerHandler) Start(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	server := &domain.Server{Id: id}
	h.service.Start(server)

	w.WriteHeader(http.StatusNoContent)
}

func (h *ServerHandler) Stop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	server := &domain.Server{Id: id}
	h.service.Stop(server)

	w.WriteHeader(http.StatusNoContent)
}

// func (h *ServerHandler) Stats(w http.ResponseWriter, r *http.Request) {
// 	c := Upgrade(w, r)
// 	defer c.Close()

// 	id := r.PathValue("id")
// 	server := &domain.Server{Id: id}
// 	stats := h.service.Stats(server)

// 	ctx, cancel := context.WithCancel(r.Context())
// 	defer cancel()

// 	go func() {
// 		for {
// 			var buf []byte
// 			_, err := c.Read(buf)
// 			if err != nil {
// 				if err == io.EOF {
// 					break
// 				}
// 			}
// 		}
// 		cancel()
// 	}()

// 	go func() {
// 		var previousStat domain.ContainerStat
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case stat, ok := <-stats:
// 				if !ok {
// 					return
// 				}

// 				if previousStat.Cpu == stat.Cpu && previousStat.Memory == stat.Memory {
// 					continue
// 				}
// 				previousStat = stat

// 				output := types.ServerStatResponse{
// 					Cpu:    stat.Cpu,
// 					Memory: stat.Memory,
// 				}

// 				data, err := json.Marshal(output)
// 				if err != nil {
// 					panic(err)
// 				}

// 				_, err = c.Write(data)
// 				if err != nil {
// 					panic(err)
// 				}
// 			}
// 		}
// 	}()

// 	<-ctx.Done()
// }

func (h *ServerHandler) Attach(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	server := &domain.Server{Id: id}

	if !h.service.Running(server) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c := Upgrade(w, r)

	ctx, cancel := context.WithCancel(r.Context())
	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	closeAll := func() {
		stdinWriter.Close()
		stdoutReader.Close()
		stdoutWriter.Close()
		c.Close()
	}

	go func() {
		defer cancel()
		io.Copy(stdinWriter, c)
	}()
	go func() {
		defer cancel()
		io.Copy(c, stdoutReader)
	}()

	go func() {
		<-ctx.Done()
		closeAll()
	}()

	done := make(chan struct{})
	err := h.service.Attach(server, stdinReader, stdoutWriter, stdoutWriter, done)
	if err != nil {
		cancel()
	}

	select {
	case <-ctx.Done():
	case <-done:
		cancel()
	}
}

func (h *ServerHandler) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("POST /servers", h.Create)
	m.HandleFunc("GET /servers", h.Read)
	m.HandleFunc("GET /servers/{id}", h.ReadOne)
	m.HandleFunc("DELETE /servers/{id}", h.Delete)

	m.HandleFunc("POST /servers/{id}/start", h.Start)
	m.HandleFunc("POST /servers/{id}/stop", h.Stop)

	// m.HandleFunc("GET /servers/{id}/stats", h.Stats)
	m.HandleFunc("GET /servers/{id}/attach", h.Attach)
}
