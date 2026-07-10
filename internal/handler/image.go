package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ayeama/panel/internal/runtime"
)

type ImageHandler struct {
	runtime runtime.Runtime
}

func NewImageHandler(runtime runtime.Runtime) ImageHandler {
	return ImageHandler{runtime}
}

func (h *ImageHandler) RegisterHandlers(mux *http.ServeMux) {
	mux.HandleFunc("GET /images", h.handleImageReadMany)
	mux.HandleFunc("GET /images/{id}", h.handleImageRead)
}

func (h *ImageHandler) handleImageReadMany(w http.ResponseWriter, r *http.Request) {
	images, err := h.runtime.ImageReadMany()
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(images); err != nil {
		log.Fatal(err)
	}
}

func (h *ImageHandler) handleImageRead(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	image, err := h.runtime.ImageRead(id)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Add("Content-Type", "application/json")

	if err = json.NewEncoder(w).Encode(image); err != nil {
		log.Fatal(err)
	}
}
