package config

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

func ConnectAndMigrate(dbPath string, migrationsDir string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	sqlBytes, err := os.ReadFile("migrations/001_init.sql")
	if err != nil {
		panic(err)
	}

	_, err = db.Exec(string(sqlBytes))
	if err != nil {
		panic(err)
	}

	fmt.Println("Migration executed")

	return db, nil
}
