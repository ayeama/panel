package handler

import (
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

func (h *ServerHandler) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("POST /servers", h.Create)
	m.HandleFunc("GET /servers", h.Read)
	m.HandleFunc("GET /servers/{id}", h.ReadOne)
	m.HandleFunc("DELETE /servers/{id}", h.Delete)
	m.HandleFunc("POST /servers/{id}/start", h.Start)
	m.HandleFunc("POST /servers/{id}/stop", h.Stop)
}
