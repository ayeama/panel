package repository

import (
	"database/sql"
	"errors"

	"github.com/ayeama/panel/api/internal/domain"
)

type SidecarRepository struct {
	db *sql.DB
}

func NewSidecarRepository(db *sql.DB) *SidecarRepository {
	return &SidecarRepository{db: db}
}

func (r *SidecarRepository) Create(id string, container_id string, server_id string) {
	_, err := r.db.Exec("INSERT INTO sidecars (id, container_id, server_id) VALUES (?, ?, ?)", id, container_id, server_id)
	if err != nil {
		panic(err)
	}
}

func (r *SidecarRepository) Read(p domain.Pagination) domain.PaginationResponse[domain.Sidecar] {
	rows, err := r.db.Query("SELECT id, container_id, server_id FROM sidecars LIMIT ? OFFSET ?", p.Limit, p.Offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	sidecars := domain.PaginationResponse[domain.Sidecar]{
		Limit:  p.Limit,
		Offset: p.Offset,
		Items:  make([]domain.Sidecar, 0),
	}

	for rows.Next() {
		var sidecar domain.Sidecar
		err = rows.Scan(&sidecar.Id, &sidecar.ContainerId, &sidecar.ServerId)
		if err != nil {
			panic(err)
		}

		sidecars.Items = append(sidecars.Items, sidecar)
	}

	err = r.db.QueryRow("SELECT COUNT(*) FROM sidecars").Scan(&sidecars.Total)
	if err != nil {
		panic(err)
	}

	return sidecars
}

func (r *SidecarRepository) ReadByServerId(server_id string) []domain.Sidecar {
	rows, err := r.db.Query("SELECT id, container_id, server_id FROM sidecars WHERE server_id=?", server_id)
	if err != nil {
		panic(err)
	}

	sidecars := make([]domain.Sidecar, 0)
	for rows.Next() {
		var sidecar domain.Sidecar
		err = rows.Scan(&sidecar.Id, &sidecar.ContainerId, &sidecar.ServerId)
		if err != nil {
			panic(err)
		}

		sidecars = append(sidecars, sidecar)
	}

	return sidecars
}

func (r *SidecarRepository) ReadOne(id string) (domain.Sidecar, error) {
	var sidecar domain.Sidecar
	if sidecar.Container == nil {
		sidecar.Container = &domain.Container{}
	}
	err := r.db.QueryRow("SELECT id, container_id, server_id FROM sidecars WHERE id=?", id).Scan(&sidecar.Id, &sidecar.ContainerId, &sidecar.ServerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return sidecar, domain.ErrNotFound
		} else {
			panic(err)
		}
	}
	return sidecar, nil
}

func (r *SidecarRepository) Delete(id string) {
	_, err := r.db.Exec("DELETE FROM sidecars WHERE id=?", id)
	if err != nil {
		panic(err)
	}
}
