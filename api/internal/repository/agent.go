package repository

import (
	"database/sql"
	"time"

	"github.com/ayeama/panel/api/internal/domain"
	"github.com/google/uuid"
)

type AgentRepository struct {
	db *sql.DB
}

func NewAgentRepository(database *sql.DB) *AgentRepository {
	return &AgentRepository{db: database}
}

func (s *AgentRepository) Create(agent *domain.Agent) {
	_, err := s.db.Exec("INSERT INTO agents (id, hostname, seen) VALUES (?, ?, ?)", agent.Id, agent.Hostname, agent.Seen.Format(time.RFC3339Nano))
	if err != nil {
		panic(err)
	}
}

func (s *AgentRepository) Read(p domain.Pagination) domain.PaginationResponse[domain.Agent] {
	rows, err := s.db.Query("SELECT id, hostname, seen FROM agents LIMIT ? OFFSET ?", p.Limit, p.Offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	agents := domain.PaginationResponse[domain.Agent]{
		Limit:  p.Limit,
		Offset: p.Offset,
		Items:  make([]domain.Agent, 0),
	}

	for rows.Next() {
		var agent domain.Agent
		var agent_seen string

		err = rows.Scan(&agent.Id, &agent.Hostname, &agent_seen)
		if err != nil {
			panic(err)
		}

		agent.Seen, err = time.Parse(time.RFC3339Nano, agent_seen)
		if err != nil {
			panic(err)
		}

		agents.Items = append(agents.Items, agent)
	}

	err = s.db.QueryRow("SELECT COUNT(*) FROM agents").Scan(&agents.Total)
	if err != nil {
		panic(err)
	}

	return agents
}

func (s *AgentRepository) ReadOne(agent *domain.Agent) {
	err := s.db.QueryRow("SELECT hostname, seen FROM agents WHERE id=?", agent.Id).Scan(&agent.Hostname, &agent.Seen)
	if err != nil {
		panic(err)
	}
}

func (s *AgentRepository) Delete(id string) {
	result, err := s.db.Exec("DELETE FROM agents WHERE id=?", id)
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

func (s *AgentRepository) Update(agent *domain.Agent, stat *domain.AgentStat) {
	var agent_seen string
	err := s.db.QueryRow("SELECT id, hostname, seen FROM agents WHERE hostname=?", agent.Hostname).Scan(&agent.Id, &agent.Hostname, &agent_seen)
	if err != nil {
		panic(err)
	}

	agent.Seen, err = time.Parse(time.RFC3339Nano, agent_seen)
	if err != nil {
		panic(err)
	}

	if agent.Id == "" {
		agent.Id = uuid.NewString()
		agent.Seen = stat.Time
		s.Create(agent)
	} else {
		_, err = s.db.Exec("UPDATE agents SET seen=? WHERE hostname=?", stat.Time.UTC().Format(time.RFC3339Nano), agent.Hostname)
		if err != nil {
			panic(err)
		}
	}
}
