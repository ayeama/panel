package service

import (
	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	repository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{repository: userRepository}
}

func (s *UserService) Create(user *domain.User) {
	user.Id = uuid.NewString()
	s.repository.Create(user)
}

func (s *UserService) Read(p domain.Pagination) domain.PaginationResponse[domain.User] {
	return s.repository.Read(p)
}

func (s *UserService) ReadOne(user *domain.User) {
	s.repository.ReadOne(user)
}

func (s *UserService) Delete(user *domain.User) {
	s.repository.Delete(user)
}
