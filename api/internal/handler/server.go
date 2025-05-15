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

func (s *ServerHandler) Create(w http.ResponseWriter, r *http.Request) {
	var serverCreate types.ServerCreateRequest
	ReadRequestJson(r.Body, &serverCreate)
	server := s.service.Create(domain.Server{Name: serverCreate.Name})
	output := types.ServerResponse{Id: server.Id, Name: server.Name, Status: server.Status}
	WriteResponseJson(w, 200, output)
}

func (s *ServerHandler) Read(w http.ResponseWriter, r *http.Request) {
	pagination := domain.NewPagination(r.URL.Query())

	// TODO add helper function for conversion
	domainServerPaginated := s.service.Read(pagination)
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

func (s *ServerHandler) ReadOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	domainServer := s.service.ReadOne(id)
	server := types.ServerResponse{Id: domainServer.Id, Name: domainServer.Name, Status: domainServer.Status}
	WriteResponseJson(w, 200, server)
}

func (s *ServerHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	var serverUpdate types.ServerUpdateRequest
	ReadRequestJson(r.Body, &serverUpdate)
	domainServer := domain.Server{Id: id, Name: serverUpdate.Name}
	s.service.Update(&domainServer)

	server := types.ServerResponse{Id: domainServer.Id, Name: domainServer.Name, Status: domainServer.Status}
	WriteResponseJson(w, 200, server)
}

func (s *ServerHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.service.Delete(id)

	w.WriteHeader(http.StatusNoContent)
}

func (s *ServerHandler) Start(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.service.Start(id)
}

func (s *ServerHandler) Stop(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	s.service.Stop(id)
}

// func (s *ServerHandler) ServerAttach(w http.ResponseWriter, r *http.Request) {
// 	conn := Upgrade(w, r)
// 	defer conn.Close()

// 	id := r.PathValue("id")

// 	done := make(chan bool, 1)
// 	stdinReader, stdinWriter := io.Pipe()
// 	stdoutReader, stdoutWriter := io.Pipe()

// 	go func() {
// 		defer stdinWriter.Close()
// 		io.Copy(stdinWriter, conn)
// 		done <- true
// 	}()
// 	go func() {
// 		defer stdoutReader.Close()
// 		io.Copy(conn, stdoutReader)
// 		done <- true
// 	}()

// 	s.service.ServerAttach(id, stdinReader, stdoutWriter, nil)
// 	<-done
// }

func (s *ServerHandler) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("POST /servers", s.Create)
	m.HandleFunc("GET /servers", s.Read)
	m.HandleFunc("GET /servers/{id}", s.ReadOne)
	m.HandleFunc("PATCH /servers/{id}", s.Update)
	m.HandleFunc("DELETE /servers/{id}", s.Delete)
	m.HandleFunc("POST /servers/{id}/start", s.Start)
	m.HandleFunc("POST /servers/{id}/stop", s.Stop)
	// m.HandleFunc("GET /servers/{id}/attach", s.ServerAttach)
}
