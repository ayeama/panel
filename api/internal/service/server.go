package service

import (
	"io"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/ayeama/panel/api/internal/runtime"
	"github.com/google/uuid"
)

type ServerService struct {
	runtime    runtime.Runtime
	repository *repository.ServerRepository
}

func NewServerService(runtime runtime.Runtime, serverRepository *repository.ServerRepository) *ServerService {
	return &ServerService{runtime: runtime, repository: serverRepository}
}

func (s *ServerService) Create(server *domain.Server) {
	server.Id = uuid.NewString()
	server.Status = "created"

	containerId := s.runtime.Create(server)

	server.Container = &domain.Container{
		Id: containerId,
	}

	s.repository.Create(server)
}

func (s *ServerService) Read(p domain.Pagination) domain.PaginationResponse[domain.Server] {
	return s.repository.Read(p)
}

func (s *ServerService) ReadOne(server *domain.Server) {
	s.repository.ReadOne(server)
}

func (s *ServerService) Update(server *domain.Server) {
	s.repository.Update(server)
}

func (s *ServerService) Delete(server *domain.Server) {
	s.repository.ReadOne(server)
	s.runtime.Delete(server.Container)
	s.repository.Delete(server)
}

func (s *ServerService) Start(server *domain.Server) {
	s.repository.ReadOne(server)
	s.runtime.Start(server.Container)
	// s.RefreshStatus(server)
}

func (s *ServerService) Stop(server *domain.Server) {
	s.repository.ReadOne(server)
	s.runtime.Stop(server.Container)
	// s.RefreshStatus(server)
}

func (s *ServerService) Stats(server *domain.Server) chan domain.ContainerStat {
	s.repository.ReadOne(server)
	return s.runtime.Stats(server.Container)
}

func (s *ServerService) Attach(server *domain.Server, stdin io.Reader, stdout io.Writer, stderr io.Writer) error {
	s.repository.ReadOne(server)
	return s.runtime.Attach(server.Container, stdin, stdout, stderr)
}

func (s *ServerService) Events() chan domain.Event {
	return s.runtime.Events()
}

// func (s *ServerService) RefreshStatus(server *domain.Server) {
// 	server.Status = s.runtime.Status(server.Container)
// 	s.repository.UpdateStatus(server)
// }
