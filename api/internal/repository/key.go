package repository

import (
	"database/sql"
	"errors"

	"github.com/ayeama/panel/api/internal/domain"
)

type KeyRepository struct {
	db *sql.DB
}

func NewKeyRepository(db *sql.DB) *KeyRepository {
	return &KeyRepository{db: db}
}

func (r *KeyRepository) Create(id string, comment string, public_key string) {
	_, err := r.db.Exec("INSERT INTO keys (id, comment, public_key) VALUES (?, ?, ?)", id, comment, public_key)
	if err != nil {
		panic(err)
	}
}

func (r *KeyRepository) Read(p domain.Pagination) domain.PaginationResponse[domain.Key] {
	rows, err := r.db.Query("SELECT id, comment, public_key FROM keys LIMIT ? OFFSET ?", p.Limit, p.Offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	keys := domain.PaginationResponse[domain.Key]{
		Limit:  p.Limit,
		Offset: p.Offset,
		Items:  make([]domain.Key, 0),
	}

	for rows.Next() {
		var key domain.Key
		err = rows.Scan(&key.Id, &key.Comment, &key.PublicKey)
		if err != nil {
			panic(err)
		}

		keys.Items = append(keys.Items, key)
	}

	err = r.db.QueryRow("SELECT COUNT(*) FROM keys").Scan(&keys.Total)
	if err != nil {
		panic(err)
	}

	return keys
}

func (r *KeyRepository) ReadOne(id string) (domain.Key, error) {
	var key domain.Key
	err := r.db.QueryRow("SELECT id, comment, public_key FROM keys WHERE id=?", id).Scan(&key.Id, &key.Comment, &key.PublicKey)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return key, domain.ErrNotFound
		} else {
			panic(err)
		}
	}
	return key, nil
}

func (r *KeyRepository) Delete(id string) {
	_, err := r.db.Exec("DELETE FROM keys WHERE id=?", id)
	if err != nil {
		panic(err)
	}
}
