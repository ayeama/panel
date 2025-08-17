package repository

import (
	"database/sql"
	"errors"

	"github.com/ayeama/panel/api/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) ReadOne(id string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, email FROM users WHERE id=?", id).Scan(&user.Id, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, domain.ErrNotFound
		} else {
			panic(err)
		}
	}
	return user, nil
}

func (r *UserRepository) ReadOneByEmail(email string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT id, email FROM users WHERE email=?", email).Scan(&user.Id, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, domain.ErrNotFound
		} else {
			panic(err)
		}
	}
	return user, nil
}

func (r *UserRepository) Create(id string, email string) {
	_, err := r.db.Exec("INSERT INTO users (id, email) VALUES (?, ?)", id, email)
	if err != nil {
		panic(err)
	}
}
