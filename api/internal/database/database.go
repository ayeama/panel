package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Db *sql.DB
}

func NewDatabase() *Database {
	db, err := sql.Open("sqlite3", "panel.db")
	if err != nil {
		panic(err)
	}
	return &Database{Db: db}
}
