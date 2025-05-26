package handler

import (
	"net/http"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/service"
	"github.com/ayeama/panel/api/pkg/api/types"
)

type ManifestHandler struct {
	service *service.ManifestService
}

func NewManifestHandler(service *service.ManifestService) *ManifestHandler {
	return &ManifestHandler{service: service}
}

func (h *ManifestHandler) Create(w http.ResponseWriter, r *http.Request) {
	var manifestCreate types.ManifestCreateRequest
	ReadRequestJson(r.Body, &manifestCreate)

	payload := domain.ManifestPayload{}
	for _, v := range manifestCreate.Variables {
		payload.Variables = append(payload.Variables, domain.ManifestVariable{
			Name:    v.Name,
			Options: v.Options,
		})
	}

	manifest := &domain.Manifest{
		Name:    manifestCreate.Name,
		Variant: &manifestCreate.Variant,
		Version: manifestCreate.Version,
		Payload: payload,
	}

	h.service.Create(manifest)

	outputVariables := make([]types.ManifestVariableResponse, 0)
	for _, v := range manifest.Payload.Variables {
		outputVariables = append(outputVariables, types.ManifestVariableResponse{
			Name:    v.Name,
			Options: v.Options,
		})
	}
	output := types.ManifestResponse{
		Id:        manifest.Id,
		Name:      manifest.Name,
		Variant:   manifest.Variant,
		Version:   manifest.Version,
		Variables: outputVariables,
	}

	WriteResponseJson(w, 200, output)
}

func (h *ManifestHandler) Read(w http.ResponseWriter, r *http.Request) {
	pagination := domain.NewPagination(r.URL.Query())

	domainManifestPagination := h.service.Read(pagination)
	manifestPaginated := types.PaginationResponse[types.ManifestResponse]{
		Limit:  domainManifestPagination.Limit,
		Offset: domainManifestPagination.Offset,
		Total:  domainManifestPagination.Total,
		Items:  make([]types.ManifestResponse, 0),
	}

	for _, m := range domainManifestPagination.Items {
		var variables []types.ManifestVariableResponse
		for _, v := range m.Payload.Variables {
			variables = append(variables, types.ManifestVariableResponse{Name: v.Name, Options: v.Options})
		}

		manifestPaginated.Items = append(manifestPaginated.Items, types.ManifestResponse{Id: m.Id, Name: m.Name, Variant: m.Variant, Version: m.Version, Variables: variables})
	}

	WriteResponseJson(w, 200, manifestPaginated)
}

func (h *ManifestHandler) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("POST /manifests", h.Create)
	m.HandleFunc("GET /manifests", h.Read)
}
