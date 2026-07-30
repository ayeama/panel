package panel

import (
	"context"
	"encoding/json"
	"net/http"
)

type Image struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ImageService service

func (s *ImageService) Read(ctx context.Context) ([]Image, error) {
	images := make([]Image, 0)

	req, err := s.client.newRequestWithContext(ctx, http.MethodGet, "/images", nil)
	if err != nil {
		return images, err
	}

	resp, err := s.client.http.Do(req)
	if err != nil {
		return images, err
	}
	defer resp.Body.Close()

	if err = json.NewDecoder(resp.Body).Decode(&images); err != nil {
		return images, err
	}

	return images, nil
}
