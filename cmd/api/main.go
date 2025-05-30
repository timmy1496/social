package main

import (
	"github.com/joho/godotenv"
	"github.com/timmy1496/social/internal/db"
	"github.com/timmy1496/social/internal/env"
	"github.com/timmy1496/social/internal/store"
	"log"
)

const version = "0.0.1"

//	@title			GopherSocial API
//	@description	API for GopherSocial, a social network for gohpers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cfg := config{
		addr:   env.GetString("ADDR", ":8077"),
		apiURL: env.GetString("EXTERNAL_URL", "localhost:8077"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
		},
		env: env.GetString("ENV", "dev"),
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
