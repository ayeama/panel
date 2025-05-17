package repository

import (
	"database/sql"

	"github.com/ayeama/panel/api/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) {
	_, err := r.db.Exec("INSERT INTO users (id, username) VALUES (?, ?)", user.Id, user.Username)
	if err != nil {
		panic(err)
	}
}

func (r *UserRepository) Read(p domain.Pagination) domain.PaginationResponse[domain.User] {
	rows, err := r.db.Query("SELECT id, username FROM users LIMIT ? OFFSET ?", p.Limit, p.Offset)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	users := domain.PaginationResponse[domain.User]{
		Limit:  p.Limit,
		Offset: p.Offset,
		Items:  make([]domain.User, 0),
	}

	for rows.Next() {
		var user domain.User

		err = rows.Scan(&user.Id, &user.Username)
		if err != nil {
			panic(err)
		}

		users.Items = append(users.Items, user)
	}

	err = r.db.QueryRow("SELECT COUNT(*) FROM users").Scan(&users.Total)
	if err != nil {
		panic(err)
	}

	return users
}

func (r *UserRepository) ReadOne(user *domain.User) {
	err := r.db.QueryRow("SELECT id, username FROM users WHERE id = ?", &user.Id).Scan(&user.Id, &user.Username)
	if err != nil {
		panic(err)
	}
}

func (r *UserRepository) Delete(user *domain.User) {
	_, err := r.db.Exec("DELETE FROM users WHERE id = ?", user.Id)
	if err != nil {
		panic(err)
	}
}
