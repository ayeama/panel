package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/ayeama/panel/api/internal/domain"
)

type ServerRepository struct {
	db *sql.DB
}

func NewServerRepository(database *sql.DB) *ServerRepository {
	return &ServerRepository{db: database}
}

func (s *ServerRepository) Create(server *domain.Server) {
	_, err := s.db.Exec("INSERT INTO servers (id, name, status) VALUES (?, ?, ?)", server.Id, server.Name, server.Status)
	if err != nil {
		panic(err)
	}
}

func (s *ServerRepository) Read(p domain.Pagination) domain.PaginationResponse[domain.Server] {
	rows, err := s.db.Query("SELECT s.id, s.name, s.status, a.id, a.hostname, c.id, c.name, c.status FROM servers s LEFT JOIN agents a ON s.agent_id = a.id LEFT JOIN containers c ON s.id = c.server_id LIMIT ? OFFSET ?", p.Limit, p.Offset)
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
		var agent_id, agent_hostname sql.NullString
		var container_id, container_name, container_status sql.NullString

		err = rows.Scan(&server.Id, &server.Name, &server.Status, &agent_id, &agent_hostname, &container_id, &container_name, &container_status)
		if err != nil {
			panic(err)
		}

		if agent_id.Valid {
			server.Agent = &domain.Agent{
				Id:       agent_id.String,
				Hostname: agent_hostname.String,
			}
		}
		if container_id.Valid {
			server.Container = &domain.Container{
				Id:     container_id.String,
				Name:   container_name.String,
				Status: container_status.String,
			}
		}

		fmt.Println(server)
		servers.Items = append(servers.Items, server)
	}

	err = s.db.QueryRow("SELECT COUNT(*) FROM servers").Scan(&servers.Total)
	if err != nil {
		panic(err)
	}

	return servers
}

func (s *ServerRepository) ReadOne(id string) domain.Server {
	server := domain.Server{}
	var container_id, container_name, container_status sql.NullString

	err := s.db.QueryRow("SELECT s.id, s.name, s.status, c.id, c.name, c.status FROM servers s LEFT JOIN containers c ON s.id = c.server_id WHERE s.id = ?", id).Scan(&server.Id, &server.Name, &server.Status, &container_id, &container_name, &container_status)
	if err != nil {
		panic(err)
	}

	if container_id.Valid {
		server.Container = &domain.Container{
			Id:     container_id.String,
			Name:   container_name.String,
			Status: container_status.String,
		}
	}

	return server
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
