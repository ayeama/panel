package service

import (
	"github.com/ayeama/panel/api/internal/broker"
	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/google/uuid"
)

type ServerService struct {
	broker     broker.Broker
	repository *repository.ServerRepository
}

func NewServerService(broker broker.Broker, repository *repository.ServerRepository) *ServerService {
	return &ServerService{broker: broker, repository: repository}
}

func (s *ServerService) Create(server domain.Server) domain.Server {
	server.Id = uuid.NewString()
	server.Status = domain.ServerStatusCreating

	s.repository.Create(&server)
	s.broker.AddEventServerCreate(domain.EventServerCreate{Id: server.Id})

	return server
}

func (s *ServerService) Read(p domain.Pagination) domain.PaginationResponse[domain.Server] {
	return s.repository.Read(p)
}

func (s *ServerService) ReadOne(id string) domain.Server {
	return s.repository.ReadOne(id)
}

// func (s *ServerService) ServerUpdate() {}

func (s *ServerService) Delete(id string) {
	// todo should actually
	// update server status to deleting
	// post event to delete
	// ep will actually remove from database
	server := s.repository.ReadOne(id)
	if server.Container != nil {
		s.broker.AddEventServerDelete(domain.EventServerDelete{Id: id, ContainerId: server.Container.Id})
	}
	s.repository.Delete(id)
}

func (s *ServerService) Start(id string) {
	server := s.repository.ReadOne(id)
	if server.Container != nil {
		s.broker.AddEventServerStart(domain.EventServerStart{Id: id, ContainerId: server.Container.Id})
	}
}

func (s *ServerService) Stop(id string) {
	server := s.repository.ReadOne(id)
	if server.Container != nil {
		s.broker.AddEventServerStop(domain.EventServerStop{Id: id, ContainerId: server.Container.Id})
	}
}

// func (s *ServerService) ServerAttach(id string, stdin io.Reader, stdout io.Writer, stderr io.Writer) {
// }
