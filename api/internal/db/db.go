package db

import (
	"fmt"
	"log"

	"github.com/RinatHar/FarmFocus/api/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib" // pgx driver
	"github.com/jmoiron/sqlx"
)

var DB *sqlx.DB

func Connect(cfg *config.Config) {
    dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)
    db, err := sqlx.Connect("pgx", dsn) // pgx driver
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    DB = db
    log.Println("Connected to Postgres!")
}
