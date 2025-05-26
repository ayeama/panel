package repository

import (
	"database/sql"
	"encoding/json"

	"github.com/ayeama/panel/api/internal/domain"
)

type ManifestRepository struct {
	db *sql.DB
}

func NewManifestRepository(db *sql.DB) *ManifestRepository {
	return &ManifestRepository{db: db}
}

func (r *ManifestRepository) Create(manifest *domain.Manifest) {
	payload, err := json.Marshal(manifest.Payload)
	if err != nil {
		panic(err)
	}

	_, err = r.db.Exec("INSERT INTO manifests (id, name, variant, version, payload) VALUES (?, ?, ?, ?, ?)", manifest.Id, manifest.Name, manifest.Variant, manifest.Version, payload)
	if err != nil {
		panic(err)
	}
}

func (r *ManifestRepository) Read(p domain.Pagination) domain.PaginationResponse[domain.Manifest] {
	rows, err := r.db.Query("SELECT id, name, variant, version, payload FROM manifests LIMIT ? OFFSET ?", p.Limit, p.Offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	manifests := domain.PaginationResponse[domain.Manifest]{
		Limit:  p.Limit,
		Offset: p.Offset,
		Items:  make([]domain.Manifest, 0),
	}

	for rows.Next() {
		var manifest domain.Manifest
		var variant sql.NullString
		var payload string

		err := rows.Scan(&manifest.Id, &manifest.Name, &variant, &manifest.Version, &payload)
		if err != nil {
			panic(err)
		}

		if variant.Valid {
			manifest.Variant = &variant.String
		} else {
			manifest.Variant = nil
		}

		err = json.Unmarshal([]byte(payload), &manifest.Payload)
		if err != nil {
			panic(err)
		}

		manifests.Items = append(manifests.Items, manifest)
	}

	err = r.db.QueryRow("SELECT COUNT(*) FROM nodes").Scan(&manifests.Total)
	if err != nil {
		panic(err)
	}

	return manifests
}
