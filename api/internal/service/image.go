package service

import (
	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
)

type ImageService struct {
	imageRepository *repository.ImageRepository
}

func NewImageService(imageRepository *repository.ImageRepository) *ImageService {
	return &ImageService{
		imageRepository: imageRepository,
	}
}

func (s *ImageService) Read(p domain.Pagination) domain.PaginationResponse[domain.Image] {
	return s.imageRepository.Read(p)
}
