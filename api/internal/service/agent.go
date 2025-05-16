package service

import (
	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
)

type AgentService struct {
	repository *repository.AgentRepository
}

func NewAgentService(repository *repository.AgentRepository) *AgentService {
	return &AgentService{repository: repository}
}

func (s *AgentService) Read(p domain.Pagination) domain.PaginationResponse[domain.Agent] {
	// s.repository.Create(&domain.Agent{Id: uuid.NewString(), Hostname: "neon", Seen: time.Now().UTC()})
	return s.repository.Read(p)
}

func (s *AgentService) ReadOne(agent *domain.Agent) {
	s.repository.ReadOne(agent)
}

func (s *AgentService) Forget() {
	// TODO
}
