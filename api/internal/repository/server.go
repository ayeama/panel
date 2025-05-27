package repository

import (
	"database/sql"
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
	tx, err := r.db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("INSERT INTO servers (id, name, status) VALUES (?, ?, ?)", server.Id, server.Name, server.Status)
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec("INSERT INTO containers (id, server_id, node_id) VALUES (?, ?, ?)", server.Container.Id, server.Id, server.Node.Id)
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}

func (r *ServerRepository) Read(p domain.Pagination) domain.PaginationResponse[domain.Server] {
	rows, err := r.db.Query("SELECT server.id, server.name, server.status, container.id, node.id, node.name, node.uri FROM servers server LEFT JOIN containers container ON container.server_id = server.id LEFT JOIN nodes node ON container.node_id = node.id LIMIT ? OFFSET ?", p.Limit, p.Offset)
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
		server.Node = &domain.Node{}

		err = rows.Scan(&server.Id, &server.Name, &server.Status, &server.Container.Id, &server.Node.Id, &server.Node.Name, &server.Node.Uri)
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

func (r *ServerRepository) ReadOne(server *domain.Server) {
	if server.Container == nil {
		server.Container = &domain.Container{}
	}
	if server.Node == nil {
		server.Node = &domain.Node{}
	}
	err := r.db.QueryRow("SELECT server.id, server.name, server.status, container.id, node.id, node.name, node.uri from servers server LEFT JOIN containers container on container.server_id = server.id LEFT JOIN nodes node on container.node_id = node.id WHERE server.id = ?", &server.Id).Scan(&server.Id, &server.Name, &server.Status, &server.Container.Id, &server.Node.Id, &server.Node.Name, &server.Node.Uri)
	if err != nil {
		panic(err)
	}
}

func (s *ServerRepository) Update(server *domain.Server) {
	query := "UPDATE servers SET "
	params := []interface{}{}
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
	tx, err := r.db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM servers WHERE id = ?", server.Id)
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec("DELETE FROM containers WHERE server_id = ?", server.Id)
	if err != nil {
		panic(err)
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
