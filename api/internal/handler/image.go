package handler

import (
	"net/http"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/service"
	"github.com/ayeama/panel/api/pkg/api/types"
)

type ImageHandler struct {
	service *service.ImageService
}

func NewImageHandler(service *service.ImageService) *ImageHandler {
	return &ImageHandler{service: service}
}

func (h *ImageHandler) Read(w http.ResponseWriter, r *http.Request) {
	pagination := domain.NewPagination(r.URL.Query())

	// TODO add helper function for conversion
	domainImagePaginated := h.service.Read(pagination)
	imagePaginated := types.PaginationResponse[types.ImageResponse]{
		Limit:  domainImagePaginated.Limit,
		Offset: domainImagePaginated.Offset,
		Total:  domainImagePaginated.Total,
		Items:  make([]types.ImageResponse, 0),
	}

	for _, s := range domainImagePaginated.Items {
		imagePaginated.Items = append(imagePaginated.Items, types.ImageResponse{Image: s.Image})
	}

	WriteResponseJson(w, 200, imagePaginated)
}

func (h *ImageHandler) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("GET /images", h.Read)
}
