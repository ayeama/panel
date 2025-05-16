package handler

import (
	"net/http"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/service"
	"github.com/ayeama/panel/api/pkg/api/types"
)

type AgentHandler struct {
	service service.AgentService
}

func NewAgentHandler(service *service.AgentService) *AgentHandler {
	return &AgentHandler{service: *service}
}

func (s *AgentHandler) Read(w http.ResponseWriter, r *http.Request) {
	pagination := domain.NewPagination(r.URL.Query())

	// TODO add helper function for conversion
	domainAgentPaginated := s.service.Read(pagination)
	agentPaginated := types.PaginationResponse[types.AgentResponse]{
		Limit:  domainAgentPaginated.Limit,
		Offset: domainAgentPaginated.Offset,
		Total:  domainAgentPaginated.Total,
		Items:  make([]types.AgentResponse, 0),
	}

	for _, s := range domainAgentPaginated.Items {
		agentPaginated.Items = append(agentPaginated.Items, types.AgentResponse{Id: s.Id, Hostname: s.Hostname, Seen: s.Seen})
	}

	WriteResponseJson(w, 200, agentPaginated)
}

func (s *AgentHandler) ReadOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	domainAgent := domain.Agent{Id: id}
	s.service.ReadOne(&domainAgent)
	agent := types.AgentResponse{Id: domainAgent.Id, Hostname: domainAgent.Hostname, Seen: domainAgent.Seen}
	WriteResponseJson(w, 200, agent)
}

func (s *AgentHandler) Forget(w http.ResponseWriter, r *http.Request) {
	// TODO
}

func (s *AgentHandler) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("GET /agents", s.Read)
	m.HandleFunc("GET /agents/{id}", s.ReadOne)
	m.HandleFunc("POST /agents/{id}/forget", s.Forget)
}
