package main

import (
	"context"
	"log"

	"github.com/RinatHar/FarmFocus/api/internal/config"
	"github.com/RinatHar/FarmFocus/api/internal/handler"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg := config.LoadConfig()
	e := echo.New()

	dbpool, err := pgxpool.New(context.Background(), cfg.GetPostgresDSN())
	if err != nil {
		log.Fatal(err)
	}
	defer dbpool.Close()

	taskRepo := repository.NewTaskRepo(dbpool)
	catRepo := repository.NewCategoryRepo(dbpool)

	taskHandler := handler.NewTaskHandler(taskRepo)
	catHandler := handler.NewCategoryHandler(catRepo)

	// Task routes
	t := e.Group("/tasks")
	t.GET("", taskHandler.GetAll)
	t.GET("/:id", taskHandler.GetByID)
	t.POST("", taskHandler.Create)
	t.PUT("/:id", taskHandler.Update)
	t.DELETE("/:id", taskHandler.Delete)

	// Category routes
	c := e.Group("/categories")
	c.GET("", catHandler.GetAll)
	c.POST("", catHandler.Create)
	c.PUT("/:id", catHandler.Update)
	c.DELETE("/:id", catHandler.Delete)

	e.Logger.Fatal(e.Start(":" + cfg.Port))
}
