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
	runtime           runtime.Runtime
	serverRepository  *repository.ServerRepository
	imageRepository   *repository.ImageRepository
	sidecarRepository *repository.SidecarRepository
}

func NewServerService(runtime runtime.Runtime, serverRepository *repository.ServerRepository, imageRepository *repository.ImageRepository, sidecarRepository *repository.SidecarRepository) *ServerService {
	return &ServerService{
		runtime:           runtime,
		serverRepository:  serverRepository,
		imageRepository:   imageRepository,
		sidecarRepository: sidecarRepository,
	}
}

func (s *ServerService) Create(tag string) domain.Server {
	id := uuid.NewString()

	image, err := s.imageRepository.ReadByTag(tag)
	if err != nil {
		panic(err)
	}

	container_id := s.runtime.Create(id, image.Tag)
	container := s.runtime.Inspect(container_id)

	server := domain.Server{
		Id:          id,
		ImageId:     image.Id,
		ContainerId: container.Id,
		Image:       &image,
		Container:   &container,
	}

	s.serverRepository.Create(server.Id, server.ImageId, server.ContainerId)

	sidecar_id := uuid.NewString()
	sidecar_container_id := s.runtime.CreateSidecar(sidecar_id, "localhost/ayeama/panel/sidecar/sftp:0.0.1", container_id)
	sidecar_container := s.runtime.Inspect(sidecar_container_id)

	sidecar := domain.Sidecar{
		Id:          sidecar_id,
		ContainerId: sidecar_container_id,
		Container:   &sidecar_container,
	}
	server.Sidecars = []*domain.Sidecar{&sidecar}

	s.sidecarRepository.Create(sidecar.Id, sidecar.ContainerId, server.Id)

	return server
}

func (s *ServerService) Read(p domain.Pagination) domain.PaginationResponse[domain.Server] {
	serversPaginated := s.serverRepository.Read(p)

	for i, server := range serversPaginated.Items {
		image, err := s.imageRepository.ReadById(server.ImageId)
		if err != nil {
			panic(err)
		}
		serversPaginated.Items[i].Image = &image

		container := s.runtime.Inspect(server.ContainerId)
		serversPaginated.Items[i].Container = &container

		sidecars := s.sidecarRepository.ReadByServerId(server.Id)
		for _, sidecar := range sidecars {
			sidecar_container := s.runtime.Inspect(sidecar.ContainerId)
			sidecar.Container = &sidecar_container
			serversPaginated.Items[i].Sidecars = append(serversPaginated.Items[i].Sidecars, &sidecar)
		}
	}

	return serversPaginated
}

func (s *ServerService) ReadOne(id string) (domain.Server, error) {
	server, err := s.serverRepository.ReadOne(id)
	if err != nil {
		return server, err
	}

	// TODO break out into a func as read also uses this logic
	image, err := s.imageRepository.ReadById(server.ImageId)
	if err != nil {
		panic(err)
	}
	server.Image = &image

	container := s.runtime.Inspect(server.ContainerId)
	server.Container = &container

	sidecars := s.sidecarRepository.ReadByServerId(server.Id)
	for _, sidecar := range sidecars {
		sidecar_container := s.runtime.Inspect(sidecar.ContainerId)
		sidecar.Container = &sidecar_container
		server.Sidecars = append(server.Sidecars, &sidecar)
	}

	return server, nil
}

func (s *ServerService) Delete(id string) error {
	server, err := s.serverRepository.ReadOne(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return err
		} else {
			panic(err)
		}
	}

	sidecars := s.sidecarRepository.ReadByServerId(server.Id)
	for _, sidecar := range sidecars {
		s.runtime.Delete(sidecar.ContainerId)
		s.sidecarRepository.Delete(sidecar.Id)
	}

	s.runtime.Delete(server.ContainerId)
	s.serverRepository.Delete(server.Id)
	return nil
}

func (s *ServerService) Start(id string) error {
	server, err := s.serverRepository.ReadOne(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return err
		} else {
			panic(err)
		}
	}

	sidecars := s.sidecarRepository.ReadByServerId(server.Id)
	for _, sidecar := range sidecars {
		s.runtime.Start(sidecar.ContainerId)
		s.runtime.InjectCredentials(sidecar.ContainerId)
	}

	s.runtime.Start(server.ContainerId)
	return nil
}

func (s *ServerService) Stop(container_id string) error {
	server, err := s.serverRepository.ReadOne(container_id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return err
		} else {
			panic(err)
		}
	}

	sidecars := s.sidecarRepository.ReadByServerId(server.Id)
	for _, sidecar := range sidecars {
		s.runtime.Stop(sidecar.ContainerId)
	}

	s.runtime.Stop(server.ContainerId)
	return nil
}

func (s *ServerService) Events() error {
	// for event := range s.runtime.Events() {
	// 	switch e := event.(type) {
	// 	case domain.RuntimeEventServerStatusChanged:
	// 		// server, err := s.serverRepository.ReadOne(e.ContainerId)
	// 		// if err != nil {
	// 		// 	if errors.Is(err, domain.ErrNotFound) {
	// 		// 		continue
	// 		// 	} else {
	// 		// 		panic(err)
	// 		// 	}
	// 		// }
	// 		// s.UpdateStatus(&server)
	// 	default:
	// 	}
	// }
	return nil
}

// func (s *ServerService) Stats(server *domain.Server) chan domain.ContainerStat {
// 	s.serverRepository.ReadOne(server)
// 	return s.runtime.Stats(server.Container)
// }

func (s *ServerService) Attach(container_id string, stdin io.Reader, stdout io.Writer, stderr io.Writer, done chan struct{}) error {
	return s.runtime.Attach(container_id, stdin, stdout, stderr, done)
}
