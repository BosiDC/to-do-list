package api

import (
	"log"
	"to-do-list/internal/config"
)

func main() {
	db, err := config.ConnectAndMigrate("./data/dev.db", "./migrations")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Database ready")
}
