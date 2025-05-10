package repository

import (
	"database/sql"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/google/uuid"
)

type ServerRepository struct {
	db *sql.DB
}

func NewServerRepository(database *sql.DB) *ServerRepository {
	return &ServerRepository{db: database}
}

func (s *ServerRepository) Create(serverTodo domain.Server) domain.Server {
	id := uuid.New().String()
	// name := namesgenerator.GetRandomName(0)
	name := serverTodo.Name
	status := "created"

	// server := domain.Server{}
	// err := s.db.QueryRow("INSERT INTO servers (id, name, status) VALUES (?, ?, ?) RETURNING id, name, status", id, name, status).Scan(&server.Id, &server.Name, &server.Status)
	// if err != nil {
	// 	panic(err)
	// }

	server := domain.Server{}
	err := s.db.QueryRow("INSERT INTO servers (id, name, status, container_id) VALUES (?, ?, ?, ?) RETURNING id, name, status, container_id", id, name, status, serverTodo.Pod.Id).Scan(&server.Id, &server.Name, &server.Pod.Status, &server.Pod.Id)
	if err != nil {
		panic(err)
	}

	return server
}

func (s *ServerRepository) Read(p domain.Pagination) domain.PaginationResponse[domain.Server] {
	rows, err := s.db.Query("SELECT id, name, status, container_id, (SELECT COUNT(*) FROM servers) AS total FROM servers LIMIT ? OFFSET ?", p.Limit, p.Offset)
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
		err = rows.Scan(&server.Id, &server.Name, &server.Pod.Status, &server.Pod.Id, &servers.Total)
		if err != nil {
			panic(err)
		}
		server.Inspect() // refreshes state
		servers.Items = append(servers.Items, server)
	}

	return servers
}

func (s *ServerRepository) ReadOne(id string) domain.Server {
	server := domain.Server{}
	err := s.db.QueryRow("SELECT id, name, status, container_id FROM servers WHERE id = ?", id).Scan(&server.Id, &server.Name, &server.Pod.Status, &server.Pod.Id)
	if err != nil {
		panic(err)
	}
	return server
}

func (s *ServerRepository) Delete(id string) {
	result, err := s.db.Exec("DELETE FROM servers WHERE id = ?", id)
	if err != nil {
		panic(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(err)
	}

	if rowsAffected == 0 {
		panic("no record found")
	}
}
