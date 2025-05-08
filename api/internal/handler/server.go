package handler

import (
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

func (s *ServerHandler) ServerCreate(w http.ResponseWriter, r *http.Request) {
	var serverCreate types.ServerCreateRequest
	ReadRequestJson(r.Body, &serverCreate)
	server := s.service.ServerCreate(domain.Server{})
	output := types.ServerResponse{Id: server.Id, Name: server.Name, Status: server.Status}
	WriteResponseJson(w, 200, output)
}

func (s *ServerHandler) ServerRead(w http.ResponseWriter, r *http.Request) {
	pagination := domain.NewPagination(r.URL.Query())

	// TODO add helper function for conversion
	domainServerPaginated := s.service.ServerRead(pagination)
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

func (s *ServerHandler) ServerReadOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	domainServer := s.service.ServerReadOne(id)
	server := types.ServerResponse{Id: domainServer.Id, Name: domainServer.Name, Status: domainServer.Status}
	WriteResponseJson(w, 200, server)
}

func (s *ServerHandler) ServerUpdate(w http.ResponseWriter, r *http.Request) {}

func (s *ServerHandler) ServerDelete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.service.ServerDelete(id)

	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerHandler) ServerStart(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.service.ServerStart(id)
}

func (s *ServerHandler) ServerStop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.service.ServerStop(id)
}

func (s *ServerHandler) ServerAttach(w http.ResponseWriter, r *http.Request) {
	conn := Upgrade(w, r)
	defer conn.Close()

	id := r.PathValue("id")

	done := make(chan bool, 1)
	stdinReader, stdinWriter := io.Pipe()
	stdoutReader, stdoutWriter := io.Pipe()

	go func() {
		defer stdinWriter.Close()
		io.Copy(stdinWriter, conn)
		done <- true
	}()
	go func() {
		defer stdoutReader.Close()
		io.Copy(conn, stdoutReader)
		done <- true
	}()

	s.service.ServerAttach(id, stdinReader, stdoutWriter, nil)
	<-done
}

func (s *ServerHandler) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("POST /servers", s.ServerCreate)
	m.HandleFunc("GET /servers", s.ServerRead)
	m.HandleFunc("GET /servers/{id}", s.ServerReadOne)
	m.HandleFunc("PATCH /servers", s.ServerUpdate)
	m.HandleFunc("DELETE /servers/{id}", s.ServerDelete)
	m.HandleFunc("POST /servers/{id}/start", s.ServerStart)
	m.HandleFunc("POST /servers/{id}/stop", s.ServerStop)
	m.HandleFunc("GET /servers/{id}/attach", s.ServerAttach)
}
