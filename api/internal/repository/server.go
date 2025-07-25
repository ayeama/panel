package repository

import (
	"database/sql"
	"errors"

	"github.com/ayeama/panel/api/internal/domain"
)

type ServerRepository struct {
	db *sql.DB
}

func NewServerRepository(db *sql.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) Create(id string, image_id string, container_id string) {
	_, err := r.db.Exec("INSERT INTO servers (id, image_id, container_id) VALUES (?, ?, ?)", id, image_id, container_id)
	if err != nil {
		panic(err)
	}
}

func (r *ServerRepository) Read(p domain.Pagination) domain.PaginationResponse[domain.Server] {
	rows, err := r.db.Query("SELECT id, image_id, container_id FROM servers LIMIT ? OFFSET ?", p.Limit, p.Offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	servers := domain.PaginationResponse[domain.Server]{
		Limit:  p.Limit,
		Offset: p.Offset,
		Items:  make([]domain.Server, 0),
	}

	for rows.Next() {
		var server domain.Server
		// server.Container = &domain.Container{}

		err = rows.Scan(&server.Id, &server.ImageId, &server.ContainerId)
		if err != nil {
			panic(err)
		}

		servers.Items = append(servers.Items, server)
	}

	err = r.db.QueryRow("SELECT COUNT(*) FROM servers").Scan(&servers.Total)
	if err != nil {
		panic(err)
	}

	return servers
}

func (r *ServerRepository) ReadOne(id string) (domain.Server, error) {
	var server domain.Server
	if server.Container == nil {
		server.Container = &domain.Container{}
	}
	err := r.db.QueryRow("SELECT id, image_id, container_id FROM servers WHERE id=?", id).Scan(&server.Id, &server.ImageId, &server.ContainerId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return server, domain.ErrNotFound
		} else {
			panic(err)
		}
	}
	return server, nil
}

func (s *ServerRepository) Update(server *domain.Server) {
	// query := "UPDATE servers SET "
	// params := []any{}
	// sets := []string{}

	// if server.Name != "" {
	// 	params = append(params, server.Name)
	// 	sets = append(sets, "name=?")
	// }
	// if server.Status != "" {
	// 	params = append(params, server.Status)
	// 	sets = append(sets, "status=?")
	// }

	// query += strings.Join(sets, ", ")
	// query += " WHERE id=? RETURNING name, status"
	// params = append(params, server.Id)

	// err := s.db.QueryRow(query, params...).Scan(&server.Name, &server.Status)
	// if err != nil {
	// 	panic(err)
	// }
}

func (r *ServerRepository) UpdateStatus(server *domain.Server) {
	// _, err := r.db.Exec("UPDATE servers SET status=? WHERE id=?", server.Status, server.Id)
	// if err != nil {
	// 	panic(err)
	// }
}

func (r *ServerRepository) Delete(id string) {
	_, err := r.db.Exec("DELETE FROM servers WHERE id=?", id)
	if err != nil {
		panic(err)
	}
}
