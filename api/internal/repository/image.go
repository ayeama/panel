package repository

import (
	"database/sql"
	"errors"

	"github.com/ayeama/panel/api/internal/domain"
)

type ImageRepository struct {
	db *sql.DB
}

func NewImageRepository(db *sql.DB) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) Read(p domain.Pagination) domain.PaginationResponse[domain.Image] {
	rows, err := r.db.Query("SELECT id, tag FROM images LIMIT ? OFFSET ?", p.Limit, p.Offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	images := domain.PaginationResponse[domain.Image]{
		Limit:  p.Limit,
		Offset: p.Offset,
		Items:  make([]domain.Image, 0),
	}

	for rows.Next() {
		var image domain.Image

		err = rows.Scan(&image.Id, &image.Tag)
		if err != nil {
			panic(err)
		}

		images.Items = append(images.Items, image)
	}

	err = r.db.QueryRow("SELECT COUNT(*) FROM images").Scan(&images.Total)
	if err != nil {
		panic(err)
	}

	return images
}

func (r *ImageRepository) ReadById(id string) (domain.Image, error) {
	var image domain.Image
	err := r.db.QueryRow("SELECT id, tag FROM images WHERE id=?", id).Scan(&image.Id, &image.Tag)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return image, err
		} else {
			panic(err)
		}
	}
	return image, nil
}

func (r *ImageRepository) ReadByTag(tag string) (domain.Image, error) {
	var image domain.Image
	err := r.db.QueryRow("SELECT id, tag FROM images WHERE tag=?", tag).Scan(&image.Id, &image.Tag)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return image, err
		} else {
			panic(err)
		}
	}
	return image, nil
}
