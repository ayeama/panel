package handler

import (
	"net/http"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/service"
	"github.com/ayeama/panel/api/pkg/api/types"
)

type NodeHandler struct {
	service *service.NodeService
}

func NewNodeHandler(service *service.NodeService) *NodeHandler {
	return &NodeHandler{service: service}
}

func (h *NodeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var nodeCreate types.NodeCreateRequest
	ReadRequestJson(r.Body, &nodeCreate)

	node := &domain.Node{Name: nodeCreate.Name, Uri: nodeCreate.Uri}
	h.service.Create(node)

	output := types.NodeResponse{Id: node.Id, Name: node.Name, Uri: node.Uri}
	WriteResponseJson(w, 200, output)
}

func (h *NodeHandler) Read(w http.ResponseWriter, r *http.Request) {
	pagination := domain.NewPagination(r.URL.Query())

	domainNodePagination := h.service.Read(pagination)
	nodePaginated := types.PaginationResponse[types.NodeResponse]{
		Limit:  domainNodePagination.Limit,
		Offset: domainNodePagination.Offset,
		Total:  domainNodePagination.Total,
		Items:  make([]types.NodeResponse, 0),
	}

	for _, n := range domainNodePagination.Items {
		nodePaginated.Items = append(nodePaginated.Items, types.NodeResponse{Id: n.Id, Name: n.Name, Uri: n.Uri})
	}

	WriteResponseJson(w, 200, nodePaginated)
}

func (h *NodeHandler) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("POST /nodes", h.Create)
	m.HandleFunc("GET /nodes", h.Read)
}
