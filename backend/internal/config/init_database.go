package config

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectAndMigrate(dbPath string, migrationsDir string) (*sql.DB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0o755); err != nil {
		return nil, err
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	migrationPath := filepath.Join(migrationsDir, "001_create_todos.sql")
	sqlBytes, err := os.ReadFile(migrationPath)
	if err != nil {
		db.Close()
		return nil, err
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		db.Close()
		return nil, err
	}

	fmt.Printf("Migration executed: %s\n", migrationPath)
	return db, nil
}
