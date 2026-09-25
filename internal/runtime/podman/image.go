package podman

import (
	"fmt"

	"github.com/ayeama/panel/internal/runtime"
	"github.com/ayeama/panel/internal/types"
	"go.podman.io/podman/v6/pkg/bindings/images"
)

func (r *Runtime) ImageReadMany() ([]types.Image, error) {
	filters := map[string][]string{"label": {imageLabelID}}
	imageListOptions := &images.ListOptions{}
	imageListOptions.WithAll(false).WithFilters(filters)

	imageList, err := images.List(*r.ctx, imageListOptions)
	if err != nil {
		return []types.Image{}, &runtime.Error{Op: "read", Resource: "image", ID: "", Err: fmt.Errorf("%w: %w", runtime.ErrInternal, err)}
	}

	images := make([]types.Image, 0)
	for _, image := range imageList {
		id := image.Labels[imageLabelID]

		if id == "" {
			continue
		}

		if len(image.Names) < 1 {
			continue
		}

		images = append(images, types.Image{
			ID:   image.Labels[imageLabelID],
			Name: image.Names[0],
		})
	}
	return images, nil
}
