package repository

import (
	"database/sql"

	"github.com/ayeama/panel/api/internal/domain"
)

type NodeRepository struct {
	db *sql.DB
}

func NewNodeRepository(db *sql.DB) *NodeRepository {
	return &NodeRepository{db: db}
}

func (r *NodeRepository) Create(node *domain.Node) {
	_, err := r.db.Exec("INSERT INTO nodes (id, name, uri) VALUES (?, ?, ?)", node.Id, node.Name, node.Uri)
	if err != nil {
		panic(err)
	}
}

func (r *NodeRepository) Read(p domain.Pagination) domain.PaginationResponse[domain.Node] {
	rows, err := r.db.Query("SELECT id, name, uri FROM nodes LIMIT ? OFFSET ?", p.Limit, p.Offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	nodes := domain.PaginationResponse[domain.Node]{
		Limit:  p.Limit,
		Offset: p.Offset,
		Items:  make([]domain.Node, 0),
	}

	for rows.Next() {
		var node domain.Node

		err := rows.Scan(&node.Id, &node.Name, &node.Uri)
		if err != nil {
			panic(err)
		}

		nodes.Items = append(nodes.Items, node)
	}

	err = r.db.QueryRow("SELECT COUNT(*) FROM nodes").Scan(&nodes.Total)
	if err != nil {
		panic(err)
	}

	return nodes
}

func (r *NodeRepository) ReadByName(name string) *domain.Node {
	var node domain.Node
	err := r.db.QueryRow("SELECT id, name, uri FROM nodes WHERE name = ?", name).Scan(&node.Id, &node.Name, &node.Uri)
	if err != nil {
		panic(err)
	}
	return &node
}
