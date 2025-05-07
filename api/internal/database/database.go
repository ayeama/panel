package database

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func initialize(database *sql.DB) error {
	_, err := database.Exec("CREATE TABLE IF NOT EXISTS migration_version (version VARCHAR(256) NOT NULL);")
	if err != nil {
		return err
	}
	return nil
}

func migrate(database *sql.DB) error {
	files, err := filepath.Glob(filepath.Join("migrations/*.sql"))
	if err != nil {
		return err
	}

	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			// TODO better error handling
			return err
		}

		// TODO parse migration file and apply logic

		_, err = database.Exec(string(content[:]))
		if err != nil {
			return err
		}

		fileBasename := filepath.Base(file)
		fileExtension := filepath.Ext(fileBasename)
		fileName := strings.TrimSuffix(fileBasename, fileExtension)
		_, err = database.Exec("INSERT INTO migration_version (version) VALUES (?);", fileName)
		if err != nil {
			return err
		}
	}

	return nil
}

func Connect() (*sql.DB, error) {
	database, err := sql.Open("sqlite3", "panel.db")
	if err != nil {
		return database, err
	}

	err = initialize(database)
	if err != nil {
		return database, err
	}

	err = migrate(database)
	if err != nil {
		return database, err
	}

	return database, err
}
