package handler

import (
	"bufio"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"

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

	output := types.ServerResponse{Id: server.Id, Name: server.Container.Name, Status: server.Container.Status, Addresses: server.Container.Ports}
	for _, sidecar := range server.Sidecars {
		output.SidecarAddresses = append(output.SidecarAddresses, sidecar.Container.Ports...)
	}
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

	for i, server := range domainServerPaginated.Items {
		serverPaginated.Items = append(serverPaginated.Items, types.ServerResponse{Id: server.Id, Name: server.Container.Name, Status: server.Container.Status, Addresses: server.Container.Ports})

		for _, sidecar := range server.Sidecars {
			serverPaginated.Items[i].SidecarAddresses = append(serverPaginated.Items[i].SidecarAddresses, sidecar.Container.Ports...)
		}
	}

	WriteResponseJson(w, 200, serverPaginated)
}

func (h *ServerHandler) ReadOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	server, err := h.service.ReadOne(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			panic(err)
		}
	}

	output := types.ServerResponse{Id: server.Id, Name: server.Container.Name, Status: server.Container.Status, Addresses: server.Container.Ports}
	for _, sidecar := range server.Sidecars {
		output.SidecarAddresses = append(output.SidecarAddresses, sidecar.Container.Ports...)
	}

	WriteResponseJson(w, 200, output)
}

func (h *ServerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.Delete(id)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ServerHandler) Start(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.Start(id)
	if err != nil {
		panic(err)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ServerHandler) Stop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	err := h.service.Stop(id)
	if err != nil {
		panic(err)
	}

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
	server, err := h.service.ReadOne(id)
	if err != nil {
		panic(err)
	}

	if server.Container.Status != "running" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	c := Upgrade(w, r)

	ctx, cancel := context.WithCancel(r.Context())
	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	stdinLogReader, stdinLogWriter := io.Pipe()
	stdin := io.TeeReader(stdinReader, stdinLogWriter)

	closeAll := func() {
		stdinWriter.Close()
		stdoutReader.Close()
		stdoutWriter.Close()

		stdinLogWriter.Close()

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

	buf := bufio.NewScanner(stdinLogReader)
	go func() {
		for buf.Scan() {
			slog.Warn(
				"terminal stdin",
				slog.String("server_id", id),
				slog.String("stdin", buf.Text()),
			)
		}
		if err := buf.Err(); err != nil {
			slog.Error("stdin log scanner")
		}
	}()

	go func() {
		<-ctx.Done()
		closeAll()
	}()

	done := make(chan struct{})
	err = h.service.Attach(server.ContainerId, stdin, stdoutWriter, stdoutWriter, done)
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
