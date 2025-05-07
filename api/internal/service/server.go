package service

import (
	"io"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
)

type ServerService struct {
	repository *repository.ServerRepository
}

func NewServerService(repository *repository.ServerRepository) *ServerService {
	return &ServerService{repository: repository}
}

func (s *ServerService) ServerCreate(server domain.Server) domain.Server {
	server.Create() // Created the actual container
	domainServer := s.repository.Create(server)
	return domainServer
}

func (s *ServerService) ServerRead(p domain.Pagination) domain.PaginationResponse[domain.Server] {
	return s.repository.Read(p)
}

func (s *ServerService) ServerReadOne(id string) domain.Server {
	return s.repository.ReadOne(id)
}

func (s *ServerService) ServerUpdate() {}

func (s *ServerService) ServerDelete(id string) {
	s.repository.Delete(id)
}

func (s *ServerService) ServerAttach(id string, stdin io.Reader, stdout io.Writer, stderr io.Writer) {
	domainServer := s.repository.ReadOne(id)
	domainServer.Attach(stdin, stdout, stderr)
}
