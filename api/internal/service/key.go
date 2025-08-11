package service

import (
	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/google/uuid"
)

type KeyService struct {
	keyRepository *repository.KeyRepository
}

func NewKeyService(keyRepository *repository.KeyRepository) *KeyService {
	return &KeyService{
		keyRepository: keyRepository,
	}
}

func (s *KeyService) Create(comment string, public_key string) domain.Key {
	id := uuid.NewString()
	s.keyRepository.Create(id, comment, public_key)
	key := domain.Key{
		Id:        id,
		Comment:   comment,
		PublicKey: public_key,
	}
	return key
}

func (s *KeyService) Read(p domain.Pagination) domain.PaginationResponse[domain.Key] {
	return s.keyRepository.Read(p)
}

func (s *KeyService) ReadOne(id string) (domain.Key, error) {
	return s.keyRepository.ReadOne(id)
}

func (s *KeyService) Delete(id string) {
	s.keyRepository.Delete(id)
}
