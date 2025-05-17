package handler

import (
	"net/http"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/service"
	"github.com/ayeama/panel/api/pkg/api/types"
)

type UserHandler struct {
	service service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: *service}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	var userCreate types.UserCreateRequest
	ReadRequestJson(r.Body, &userCreate)

	user := &domain.User{Username: userCreate.Username}
	h.service.Create(user)

	output := types.UserResponse{Id: user.Id, Username: user.Username}
	WriteResponseJson(w, 200, output)
}

func (h *UserHandler) Read(w http.ResponseWriter, r *http.Request) {
	pagination := domain.NewPagination(r.URL.Query())

	// TODO add helper function for conversion
	domainUserPaginated := h.service.Read(pagination)
	userPaginated := types.PaginationResponse[types.UserResponse]{
		Limit:  domainUserPaginated.Limit,
		Offset: domainUserPaginated.Offset,
		Total:  domainUserPaginated.Total,
		Items:  make([]types.UserResponse, 0),
	}

	for _, s := range domainUserPaginated.Items {
		userPaginated.Items = append(userPaginated.Items, types.UserResponse{Id: s.Id, Username: s.Username})
	}

	WriteResponseJson(w, 200, userPaginated)
}

func (h *UserHandler) ReadOne(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user := &domain.User{Id: id}
	h.service.ReadOne(user)

	output := types.UserResponse{Id: user.Id, Username: user.Username}
	WriteResponseJson(w, 200, output)
}

func (h *UserHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")

	user := &domain.User{Id: id}
	h.service.Delete(user)

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) RegisterHandlers(m *http.ServeMux) {
	m.HandleFunc("POST /users", h.Create)
	m.HandleFunc("GET /users", h.Read)
	m.HandleFunc("GET /users/{id}", h.ReadOne)
	m.HandleFunc("DELETE /users/{id}", h.Delete)
}
