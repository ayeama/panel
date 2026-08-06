package handler

import (
	"errors"
	"log"
	"net/http"

	"github.com/ayeama/panel/internal/runtime"
)

func handleError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, runtime.ErrBadRequest):
		http.Error(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, runtime.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, runtime.ErrExists):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, runtime.ErrInternal):
		log.Println("error", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	default:
		log.Println("error", err.Error())
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
