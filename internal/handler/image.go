package handler

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/ayeama/panel/internal/types"
	"go.podman.io/podman/v6/pkg/bindings/images"
)

type ImageHandler struct {
	runtime *context.Context
}

func NewImageHandler(runtime *context.Context) ImageHandler {
	return ImageHandler{runtime}
}

func (h *ImageHandler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /images", h.handleImageReadMany)
	mux.HandleFunc("GET /images/{id}", h.handleImageRead)
}

func (h *ImageHandler) handleImageReadMany(w http.ResponseWriter, r *http.Request) {
	filters := map[string][]string{"label": {"com.github.ayeama.panel.image.id"}}
	imageListOptions := &images.ListOptions{}
	imageListOptions.WithAll(false).WithFilters(filters)

	imageList, err := images.List(*h.runtime, imageListOptions)
	if err != nil {
		log.Fatal(err)
	}

	imageResponse := make([]types.Image, 0)
	for _, image := range imageList {
		id := image.Labels["com.github.ayeama.panel.image.id"]

		if id == "" {
			continue
		}

		if len(image.Names) < 1 {
			continue
		}

		imageResponse = append(imageResponse, types.Image{
			ID:   image.Labels["com.github.ayeama.panel.image.id"],
			Name: image.Names[0],
		})
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(imageResponse); err != nil {
		log.Fatal(err)
	}
}

func (h *ImageHandler) handleImageRead(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
}
