package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/RinatHar/FarmFocus/api/internal/config"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

func main() {
    // Загружаем конфиг
    cfg := config.LoadConfig()

    // Проверка на пустые значения
    if cfg.DBHost == "" || cfg.DBPort == "" || cfg.DBUser == "" || cfg.DBPass == "" || cfg.DBName == "" {
        log.Fatal("One or more DB connection variables are missing")
    }

    // Формируем строку подключения к Postgres
    dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBPort, cfg.DBName)

    db, err := sql.Open("pgx", dbURL)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }
    defer db.Close()

    // Путь к папке с миграциями
    migrationsDir := "./migrations"

    // Применяем все миграции
    if err := goose.Up(db, migrationsDir); err != nil {
        log.Fatalf("Failed to apply migrations: %v", err)
    }

    log.Println("Database migrated successfully!")
}
