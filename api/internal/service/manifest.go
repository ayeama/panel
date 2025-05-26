package service

import (
	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/google/uuid"
)

type ManifestService struct {
	repository *repository.ManifestRepository
}

func NewManifestService(repository *repository.ManifestRepository) *ManifestService {
	return &ManifestService{repository: repository}
}

func (s *ManifestService) Create(manifest *domain.Manifest) {
	manifest.Id = uuid.NewString()
	s.repository.Create(manifest)
}

func (s *ManifestService) Read(p domain.Pagination) domain.PaginationResponse[domain.Manifest] {
	return s.repository.Read(p)
}
