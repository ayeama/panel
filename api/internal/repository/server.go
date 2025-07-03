package repository

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/ayeama/panel/api/internal/domain"
)

type ServerRepository struct {
	db *sql.DB
}

func NewServerRepository(db *sql.DB) *ServerRepository {
	return &ServerRepository{db: db}
}

func (r *ServerRepository) Create(server *domain.Server) {
	_, err := r.db.Exec("INSERT INTO servers (id, name, status, container_id, container_port) VALUES (?, ?, ?, ?, ?)", server.Id, server.Name, server.Status, server.Container.Id, server.Container.Port)
	if err != nil {
		panic(err)
	}
}

func (r *ServerRepository) Read(p domain.Pagination) domain.PaginationResponse[domain.Server] {
	rows, err := r.db.Query("SELECT server.id, server.name, server.status, container_id, container_port FROM servers server LIMIT ? OFFSET ?", p.Limit, p.Offset)
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
		server.Container = &domain.Container{}

		err = rows.Scan(&server.Id, &server.Name, &server.Status, &server.Container.Id, &server.Container.Port)
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

func (r *ServerRepository) ReadOne(server *domain.Server) error {
	if server.Container == nil {
		server.Container = &domain.Container{}
	}
	err := r.db.QueryRow("SELECT server.id, server.name, server.status, container_id, container_port FROM servers server WHERE server.id = ?", &server.Id).Scan(&server.Id, &server.Name, &server.Status, &server.Container.Id, &server.Container.Port)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.ErrNotFound
		} else {
			panic(err)
		}
	}
	return nil
}

func (s *ServerRepository) Update(server *domain.Server) {
	query := "UPDATE servers SET "
	params := []any{}
	sets := []string{}

	if server.Name != "" {
		params = append(params, server.Name)
		sets = append(sets, "name=?")
	}
	if server.Status != "" {
		params = append(params, server.Status)
		sets = append(sets, "status=?")
	}

	query += strings.Join(sets, ", ")
	query += " WHERE id=? RETURNING name, status"
	params = append(params, server.Id)

	err := s.db.QueryRow(query, params...).Scan(&server.Name, &server.Status)
	if err != nil {
		panic(err)
	}
}

func (r *ServerRepository) UpdateStatus(server *domain.Server) {
	_, err := r.db.Exec("UPDATE servers SET status=? WHERE id=?", server.Status, server.Id)
	if err != nil {
		panic(err)
	}
}

func (r *ServerRepository) Delete(server *domain.Server) {
	_, err := r.db.Exec("DELETE FROM servers WHERE id = ?", server.Id)
	if err != nil {
		panic(err)
	}
}
