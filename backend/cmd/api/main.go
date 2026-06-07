package main

import (
	"log"
	"net/http"
	"path/filepath"
	"runtime"

	"to-do-list/internal/config"
	"to-do-list/internal/handler"
	"to-do-list/internal/repository"
	"to-do-list/internal/service"
)

func main() {
	_, filename, _, _ := runtime.Caller(0)
	base := filepath.Dir(filename)
	migrationsPath := filepath.Join(base, "../../internal/migrations")
	dbPath := filepath.Join(base, "../../data/dev.db")

	db, err := config.ConnectAndMigrate(dbPath, migrationsPath)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	repo := repository.NewTodosRepository(db)
	svc := service.NewTodosService(repo)
	h := handler.NewTodosHandler(svc)

	mux := http.NewServeMux()
	mux.Handle("/todos", h)
	mux.Handle("/todos/", h)

	log.Println("Starting server on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
