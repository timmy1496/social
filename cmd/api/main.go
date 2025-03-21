package main

import (
	"github.com/joho/godotenv"
	"github.com/timmy1496/social/internal/db"
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
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
	}

	dbConn, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)

	if err != nil {
		log.Panic(err)
	}

	defer dbConn.Close()
	log.Println("database connection pool established")

	storage := store.NewStorage(dbConn)

	app := &application{
		config:  cfg,
		storage: storage,
	}

	mux := app.mount()

	log.Fatal(app.run(mux))
}
