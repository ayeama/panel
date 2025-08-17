package service

import (
	"time"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/google/uuid"
)

type SessionService struct {
	SessionRepository *repository.SessionRepository
}

func NewSessionService(SessionRepository *repository.SessionRepository) *SessionService {
	return &SessionService{
		SessionRepository: SessionRepository,
	}
}

func (s *SessionService) Create(user_id string) domain.Session {
	id := uuid.NewString()
	expires := time.Now().UTC().Add(time.Hour * 720)
	Session := domain.Session{
		Id:      id,
		UserId:  user_id,
		Expires: expires,
	}
	s.SessionRepository.Create(Session.Id, Session.UserId, Session.Expires)
	return Session
}

func (s *SessionService) ReadOne(id string) (domain.Session, error) {
	Session, err := s.SessionRepository.ReadOne(id)
	if err != nil {
		return Session, err
	}
	return Session, nil
}

func (s *SessionService) Delete(id string) error {
	s.SessionRepository.Delete(id)
	return nil
}
