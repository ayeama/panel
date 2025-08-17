package repository

import (
	"database/sql"
	"errors"
	"time"

	"github.com/ayeama/panel/api/internal/domain"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

func (r *SessionRepository) Create(id string, user_id string, expires time.Time) {
	_, err := r.db.Exec("INSERT INTO sessions (id, user_id, expires) VALUES (?, ?, ?)", id, user_id, expires.UTC().Format(time.RFC3339))
	if err != nil {
		panic(err)
	}
}

func (r *SessionRepository) ReadOne(id string) (domain.Session, error) {
	var Session domain.Session
	var expires string
	err := r.db.QueryRow("SELECT id, user_id, expires FROM sessions WHERE id=?", id).Scan(&Session.Id, &Session.UserId, &expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Session, domain.ErrNotFound
		} else {
			panic(err)
		}
	}
	Session.Expires, err = time.Parse(time.RFC3339, expires)
	if err != nil {
		panic(err)
	}
	return Session, nil
}

func (r *SessionRepository) Delete(id string) {
	_, err := r.db.Exec("DELETE FROM sessions WHERE id=?", id)
	if err != nil {
		panic(err)
	}
}
