package service

import (
	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/google/uuid"
)

type NodeService struct {
	repository *repository.NodeRepository
}

func NewNodeService(repository *repository.NodeRepository) *NodeService {
	return &NodeService{repository: repository}
}

func (s *NodeService) Create(node *domain.Node) {
	node.Id = uuid.NewString()
	s.repository.Create(node)
}

func (s *NodeService) Read(p domain.Pagination) domain.PaginationResponse[domain.Node] {
	return s.repository.Read(p)
}
