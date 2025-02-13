package main

import (
	"github.com/joho/godotenv"
	"github.com/timmy1496/social/internal/env"
	"github.com/timmy1496/social/internal/store"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr: env.GetString("ADDR", ":8077"),
	}

	storage := store.NewStorage(nil)

	app := &application{
		config:  cfg,
		storage: storage,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
