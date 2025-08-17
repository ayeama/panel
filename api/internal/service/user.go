package service

import (
	"errors"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/google/uuid"
)

type UserService struct {
	userRepository *repository.UserRepository
}

func NewUserService(userRepository *repository.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepository,
	}
}

func (s *UserService) ReadOne(id string) domain.User {
	var user domain.User
	user, err := s.userRepository.ReadOne(id)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return user
		} else {
			panic(err)
		}
	}
	return user
}

func (s *UserService) Login(email string) domain.User {
	user, err := s.userRepository.ReadOneByEmail(email)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			id := uuid.NewString()

			user = domain.User{
				Id:    id,
				Email: email,
			}
			s.userRepository.Create(user.Id, user.Email)
		} else {
			panic(err)
		}
	}

	return user
}
