package service

import (
	"github.com/ayeama/panel/api/internal/domain"
	"github.com/ayeama/panel/api/internal/repository"
	"github.com/ayeama/panel/api/internal/runtime"
)

type ImageService struct {
	runtime         runtime.Runtime
	imageRepository *repository.ImageRepository
}

func NewImageService(runtime runtime.Runtime, imageRepository *repository.ImageRepository) *ImageService {
	return &ImageService{
		runtime:         runtime,
		imageRepository: imageRepository,
	}
}

func (s *ImageService) Read(p domain.Pagination) domain.PaginationResponse[domain.Image] {
	imagesPaginated := s.imageRepository.Read(p)

	for i, image := range imagesPaginated.Items {
		imageInspected := s.runtime.InspectImage(image.Tag)
		imagesPaginated.Items[i].Variables = imageInspected.Variables
	}

	return imagesPaginated
}
