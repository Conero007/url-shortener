package database

import (
	"database/sql"
	"encoding/json"
	"io"
	"os"
)

type Migration struct {
	Query    string `json:"query"`
	Rollback string `json:"rollback"`
}

func (m Migration) RunQuery(db *sql.DB) error {
	_, err := db.Exec(m.Query)
	return err
}

func (m Migration) RollbackQuery(db *sql.DB) error {
	_, err := db.Exec(m.Rollback)
	return err
}

type Migrations map[string]Migration

func getMigrations() (Migrations, error) {
	var m Migrations

	file, err := os.Open("database/migrations.json")
	if err != nil {
		return m, err
	}
	defer file.Close()

	fileContent, err := io.ReadAll(file)
	if err != nil {
		return m, err
	}

	if err := json.Unmarshal(fileContent, &m); err != nil {
		return m, err
	}

	return m, nil
}

func RunMigrations(db *sql.DB, dbName string) error {
	if err := createAndUseDatabase(db, dbName); err != nil {
		return err
	}

	migrations, err := getMigrations()
	if err != nil {
		return err
	}

	for _, migration := range migrations {
		if err := migration.RunQuery(db); err != nil {
			return err
		}
	}

	return err
}

func RollbackMigrations(db *sql.DB, migrationNames ...string) error {
	migrations, err := getMigrations()
	if err != nil {
		return err
	}

	for _, migrationName := range migrationNames {
		if _, ok := migrations[migrationName]; ok {
			if err = migrations[migrationName].RollbackQuery(db); err != nil {
				return err
			}
		}

	}

	return nil
}

func createAndUseDatabase(db *sql.DB, dbName string) error {
	if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS " + dbName); err != nil {
		return err
	}
	if _, err := db.Exec("USE " + dbName); err != nil {
		return err
	}
	return nil
}
