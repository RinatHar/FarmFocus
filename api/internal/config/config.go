package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    Port     string
    DBHost   string
    DBPort   string
    DBUser   string
    DBPass   string
    DBName   string
    BotToken string
}

func LoadConfig() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }

    cfg := &Config{
        Port:     os.Getenv("PORT"),
        DBHost:   os.Getenv("DB_HOST"),
        DBPort:   os.Getenv("DB_PORT"),
        DBUser:   os.Getenv("DB_USER"),
        DBPass:   os.Getenv("DB_PASSWORD"),
        DBName:   os.Getenv("DB_NAME"),
        BotToken: os.Getenv("MAX_BOT_TOKEN"),
    }

    return cfg
}

func (c *Config) GetPostgresDSN() string {
    return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
        c.DBUser,
        c.DBPass,
        c.DBHost,
        c.DBPort,
        c.DBName,
    )
}