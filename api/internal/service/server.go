package service

import (
	"io"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/ayeama/panel/api/internal/runtime"
	"github.com/google/uuid"
)

type ServerService struct {
	runtime            runtime.Runtime
	serverRepository   *repository.ServerRepository
	nodeRepository     *repository.NodeRepository
	manifestRepository *repository.ManifestRepository
}

func NewServerService(
	runtime runtime.Runtime,
	serverRepository *repository.ServerRepository,
	nodeRepository *repository.NodeRepository,
	manifestRepository *repository.ManifestRepository,
) *ServerService {
	return &ServerService{
		runtime:            runtime,
		serverRepository:   serverRepository,
		nodeRepository:     nodeRepository,
		manifestRepository: manifestRepository,
	}
}

func (s *ServerService) Create(server *domain.Server) {
	server.Id = uuid.NewString()
	server.Status = "created"

	node, err := s.runtime.RandomNode()
	if err != nil {
		panic(err)
	}
	server.Node = s.nodeRepository.ReadByName(node.Name())

	manifest := s.manifestRepository.ReadOne(server.Manifest.Id)
	server.Manifest = &manifest

	// containerId := s.runtime.Create(server)
	containerId := node.Create(server)

	server.Container = &domain.Container{
		Id: containerId,
	}

	s.serverRepository.Create(server)
}

func (s *ServerService) Read(p domain.Pagination) domain.PaginationResponse[domain.Server] {
	return s.serverRepository.Read(p)
}

func (s *ServerService) ReadOne(server *domain.Server) {
	s.serverRepository.ReadOne(server)
}

func (s *ServerService) Update(server *domain.Server) {
	s.serverRepository.Update(server)
}

func (s *ServerService) Delete(server *domain.Server) {
	s.serverRepository.ReadOne(server)

	node, err := s.runtime.Node(server.Node.Name)
	if err != nil {
		panic(err)
	}

	// s.runtime.Delete(server.Container)
	node.Delete(server.Container)
	s.serverRepository.Delete(server)
}

func (s *ServerService) Start(server *domain.Server) {
	s.serverRepository.ReadOne(server)

	node, err := s.runtime.Node(server.Node.Name)
	if err != nil {
		panic(err)
	}

	// s.runtime.Start(server.Container)
	node.Start(server.Container)
	// s.RefreshStatus(server)
}

func (s *ServerService) Stop(server *domain.Server) {
	s.serverRepository.ReadOne(server)

	node, err := s.runtime.Node(server.Node.Name)
	if err != nil {
		panic(err)
	}

	// s.runtime.Stop(server.Container)
	node.Stop(server.Container)
	// s.RefreshStatus(server)
}

func (s *ServerService) Stats(server *domain.Server) chan domain.ContainerStat {
	s.serverRepository.ReadOne(server)

	node, err := s.runtime.Node(server.Node.Name)
	if err != nil {
		panic(err)
	}

	// return s.runtime.Stats(server.Container)
	return node.Stats(server.Container)
}

func (s *ServerService) Attach(server *domain.Server, stdin io.Reader, stdout io.Writer, stderr io.Writer, done chan struct{}) error {
	s.serverRepository.ReadOne(server)

	node, err := s.runtime.Node(server.Node.Name)
	if err != nil {
		panic(err)
	}

	// return s.runtime.Attach(server.Container, stdin, stdout, stderr)
	return node.Attach(server.Container, stdin, stdout, stderr, done)
}

func (s *ServerService) Events() chan domain.Event {
	// TODO
	node, err := s.runtime.Node("rt1")
	if err != nil {
		panic(err)
	}

	// return s.runtime.Events()
	return node.Events()
}

// func (s *ServerService) RefreshStatus(server *domain.Server) {
// 	server.Status = s.runtime.Status(server.Container)
// 	s.repository.UpdateStatus(server)
// }
