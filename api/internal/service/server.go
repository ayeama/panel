package service

import (
	"errors"
	"io"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/ayeama/panel/api/internal/runtime"
	"github.com/google/uuid"
)

type ServerService struct {
	runtime          *runtime.Runtime
	serverRepository *repository.ServerRepository
	imageService     *ImageService
}

func NewServerService(runtime *runtime.Runtime, serverRepository *repository.ServerRepository, imageService *ImageService) *ServerService {
	return &ServerService{
		runtime:          runtime,
		serverRepository: serverRepository,
		imageService:     imageService,
	}
}

func (s *ServerService) Create(image string) domain.Server {
	id := uuid.NewString()
	containerId := s.runtime.Create(id, image)
	name := s.runtime.Name(containerId)
	status := s.runtime.Status(containerId)
	containerPort := s.runtime.Port(containerId)

	server := domain.Server{Id: id, Name: name, Status: status}
	server.Container = &domain.Container{Id: containerId, Port: containerPort}
	s.serverRepository.Create(&server)
	return server
}

func (s *ServerService) Read(p domain.Pagination) domain.PaginationResponse[domain.Server] {
	return s.serverRepository.Read(p)
}

func (s *ServerService) ReadOne(server *domain.Server) error {
	return s.serverRepository.ReadOne(server)
}

func (s *ServerService) Delete(server *domain.Server) error {
	err := s.serverRepository.ReadOne(server)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return err
		} else {
			panic(err)
		}
	}

	s.runtime.Delete(server.Container)
	s.serverRepository.Delete(server)
	return nil
}

func (s *ServerService) Start(server *domain.Server) error {
	err := s.serverRepository.ReadOne(server)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return err
		} else {
			panic(err)
		}
	}

	s.runtime.Start(server.Container)
	return nil
}

func (s *ServerService) Stop(server *domain.Server) error {
	err := s.serverRepository.ReadOne(server)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return err
		} else {
			panic(err)
		}
	}

	s.runtime.Stop(server.Container)
	return nil
}

func (s *ServerService) Events() error {
	for event := range s.runtime.Events() {
		switch e := event.(type) {
		case domain.RuntimeEventServerStatusChanged:
			server := domain.Server{Id: e.ServerId}
			err := s.serverRepository.ReadOne(&server)
			if err != nil {
				if errors.Is(err, domain.ErrNotFound) {
					continue
				} else {
					panic(err)
				}
			}
			s.UpdateStatus(&server)
		default:
		}
	}
	return nil
}

func (s *ServerService) UpdateStatus(server *domain.Server) {
	status := s.runtime.Status(server.Container.Id)
	server.Status = status
	s.serverRepository.UpdateStatus(server)
}

// func (s *ServerService) Stats(server *domain.Server) chan domain.ContainerStat {
// 	s.serverRepository.ReadOne(server)
// 	return s.runtime.Stats(server.Container)
// }

func (s *ServerService) Attach(server *domain.Server, stdin io.Reader, stdout io.Writer, stderr io.Writer, done chan struct{}) error {
	s.ReadOne(server)
	return s.runtime.Attach(server.Container, stdin, stdout, stderr, done)
}

func (s *ServerService) Running(server *domain.Server) bool {
	s.ReadOne(server)
	return s.runtime.Running(server.Container.Id)
}
