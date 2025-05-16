package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	Db *sql.DB
}

func NewDatabase() *Database {
	db, err := sql.Open("sqlite3", "file:panel.db?_journal=WAL&_sync=NORMAL&_txlock=immediate&_timeout=5000&_loc=auto")
	if err != nil {
		panic(err)
	}
	return &Database{Db: db}
}
