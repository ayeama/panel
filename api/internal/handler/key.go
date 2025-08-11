package handler

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/service"
	"github.com/ayeama/panel/api/pkg/api/types"
)

func hashKey(key string) string {
	h := sha256.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}

type KeyHandler struct {
	service *service.KeyService
}

func NewKeyHandler(service *service.KeyService) *KeyHandler {
	return &KeyHandler{service: service}
}

func (h *KeyHandler) Create(w http.ResponseWriter, r *http.Request) {
	var keyCreate types.KeyCreateRequest
	ReadRequestJson(r.Body, &keyCreate)

	key := h.service.Create(keyCreate.Comment, keyCreate.PublicKey)

	output := types.KeyResponse{Id: key.Id, Comment: key.Comment, HashedPublicKey: hashKey(key.PublicKey)}
	WriteResponseJson(w, 200, output)
}

func (h *KeyHandler) Read(w http.ResponseWriter, r *http.Request) {
	pagination := domain.NewPagination(r.URL.Query())

	// TODO add helper function for conversion
	domainKeyPaginated := h.service.Read(pagination)
	keyPaginated := types.PaginationResponse[types.KeyResponse]{
		Limit:  domainKeyPaginated.Limit,
		Offset: domainKeyPaginated.Offset,
		Total:  domainKeyPaginated.Total,
		Items:  make([]types.KeyResponse, 0),
	}

	for _, key := range domainKeyPaginated.Items {
		keyPaginated.Items = append(keyPaginated.Items, types.KeyResponse{Id: key.Id, Comment: key.Comment, HashedPublicKey: hashKey(key.PublicKey)})
	}

	WriteResponseJson(w, 200, keyPaginated)
}

func (h *KeyHandler) ReadOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	key, err := h.service.ReadOne(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		} else {
			panic(err)
		}
	}

	output := types.KeyResponse{Id: key.Id, Comment: key.Comment, HashedPublicKey: hashKey(key.PublicKey)}
	WriteResponseJson(w, 200, output)
}

func (h *KeyHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	h.service.Delete(id)
	w.WriteHeader(http.StatusNoContent)
}

func (h *KeyHandler) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("POST /keys", h.Create)
	m.HandleFunc("GET /keys", h.Read)
	m.HandleFunc("GET /keys/{id}", h.ReadOne)
	m.HandleFunc("DELETE /keys/{id}", h.Delete)
}
