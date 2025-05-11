package main

import (
	"github.com/timmy1496/social/internal/db"
	"github.com/timmy1496/social/internal/env"
	"github.com/timmy1496/social/internal/store"
	"log"
)

func main() {
	addr := env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable")
	conn, err := db.New(addr, 3, 3, "15m")
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	store := store.NewStorage(conn)

	db.Seed(&store)
}
