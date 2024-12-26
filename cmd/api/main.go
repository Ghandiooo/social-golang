package main

import (
	"log"
	"time"

	"github.com/ghandiooo/gosocial/internal/db"
	"github.com/ghandiooo/gosocial/internal/env"
	"github.com/ghandiooo/gosocial/internal/store"
)

func main() {
	cfg := config{
		address: env.GetString("ADDR", ":8080"),
		db: dbConfig{
			addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost:5433/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			maxIdleTime:  env.GetDuration("DB_MAX_IDLE_TIME", 10*time.Second),
		},
	}
	database, err := db.New(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	log.Println("Database connection pool established")
	storage := store.NewStorage(database)
	app := &application{
		config: cfg,
		store:  storage,
	}
	mux := app.mount()
	err = app.run(mux)
	if err != nil {
		log.Fatal(err)
	}
}
