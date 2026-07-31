package panel

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ayeama/panel/pkg/api"
)

type Image struct {
	ID   string
	Name string
}

type ImageService service

func (s *ImageService) Read(ctx context.Context) ([]Image, error) {
	var images []Image

	req, err := s.client.newRequestWithContext(ctx, http.MethodGet, "/images", nil)
	if err != nil {
		return images, err
	}

	resp, err := s.client.http.Do(req)
	if err != nil {
		return images, err
	}
	defer resp.Body.Close()

	var respImages []api.Image
	if err = json.NewDecoder(resp.Body).Decode(&respImages); err != nil {
		return images, err
	}

	images = make([]Image, len(respImages))
	for i, respImage := range respImages {
		images[i] = Image{
			ID:   respImage.ID,
			Name: respImage.Name,
		}
	}

	return images, nil
}
