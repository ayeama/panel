package handler

import (
	"encoding/json"
	"net/http"

	"github.com/ayeama/panel/internal/runtime"
	"github.com/ayeama/panel/pkg/api"
)

type ImageHandler struct {
	runtime runtime.Runtime
}

func NewImageHandler(runtime runtime.Runtime) ImageHandler {
	return ImageHandler{runtime}
}

func (h *ImageHandler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /images", h.handleImageReadMany)
}

func (h *ImageHandler) handleImageReadMany(w http.ResponseWriter, r *http.Request) {
	images, err := h.runtime.ImageReadMany()
	if err != nil {
		handleError(w, err)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	resp := make([]api.Image, len(images))
	for i, image := range images {
		resp[i] = api.Image{
			ID:   image.ID,
			Name: image.Name,
		}
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		handleError(w, err)
		return
	}
}
