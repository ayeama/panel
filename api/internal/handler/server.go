package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/service"
	"github.com/ayeama/panel/api/pkg/api/types"
)

type ServerHandler struct {
	service service.ServerService
}

func NewServerHandler(service *service.ServerService) *ServerHandler {
	return &ServerHandler{service: *service}
}

func (h *ServerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var serverCreate types.ServerCreateRequest
	ReadRequestJson(r.Body, &serverCreate)

	server := &domain.Server{Name: serverCreate.Name}
	h.service.Create(server)

	output := types.ServerResponse{Id: server.Id, Name: server.Name, Status: server.Status}
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
		serverPaginated.Items = append(serverPaginated.Items, types.ServerResponse{Id: s.Id, Name: s.Name, Status: s.Status})
	}

	WriteResponseJson(w, 200, serverPaginated)
}

func (h *ServerHandler) ReadOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	server := &domain.Server{Id: id}
	h.service.ReadOne(server)

	output := types.ServerResponse{Id: server.Id, Name: server.Name, Status: server.Status}
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

func (h *ServerHandler) Stats(w http.ResponseWriter, r *http.Request) {
	c := Upgrade(w, r)
	defer c.Close()

	id := r.PathValue("id")
	server := &domain.Server{Id: id}
	stats := h.service.Stats(server)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		for {
			var buf []byte
			_, err := c.Read(buf)
			if err != nil {
				if err == io.EOF {
					break
				}
			}
		}
		cancel()
	}()

	go func() {
		var previousStat domain.ContainerStat
		for {
			select {
			case <-ctx.Done():
				return
			case stat, ok := <-stats:
				if !ok {
					return
				}

				if previousStat.Cpu == stat.Cpu && previousStat.Memory == stat.Memory {
					continue
				}
				previousStat = stat

				output := types.ServerStatResponse{
					Cpu:    stat.Cpu,
					Memory: stat.Memory,
				}

				data, err := json.Marshal(output)
				if err != nil {
					panic(err)
				}

				_, err = c.Write(data)
				if err != nil {
					panic(err)
				}
			}
		}
	}()

	<-ctx.Done()
}

func (h *ServerHandler) Attach(w http.ResponseWriter, r *http.Request) {
	c := Upgrade(w, r)
	defer c.Close()

	id := r.PathValue("id")
	server := &domain.Server{Id: id}

	ctx, cancel := context.WithCancel(context.Background())
	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	go func() {
		defer stdinWriter.Close()
		_, err := io.Copy(stdinWriter, c)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
		}
		cancel()
	}()
	go func() {
		defer stdoutReader.Close()
		_, err := io.Copy(c, stdoutReader)
		if err != nil {
			if err != io.EOF {
				panic(err)
			}
		}
		cancel()
	}()

	h.service.Attach(server, stdinReader, stdoutWriter, stdoutWriter)
	<-ctx.Done()
}

func (h *ServerHandler) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("POST /servers", h.Create)
	m.HandleFunc("GET /servers", h.Read)
	m.HandleFunc("GET /servers/{id}", h.ReadOne)
	m.HandleFunc("DELETE /servers/{id}", h.Delete)

	m.HandleFunc("POST /servers/{id}/start", h.Start)
	m.HandleFunc("POST /servers/{id}/stop", h.Stop)

	m.HandleFunc("GET /servers/{id}/stats", h.Stats)
	m.HandleFunc("GET /servers/{id}/attach", h.Attach)
}
