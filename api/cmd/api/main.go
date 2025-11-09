package main

import (
	"github.com/RinatHar/FarmFocus/api/internal/config"
	"github.com/RinatHar/FarmFocus/api/internal/db"
	"github.com/RinatHar/FarmFocus/api/internal/handler"
	"github.com/labstack/echo/v4"
)

func main() {
    cfg := config.LoadConfig()
    db.Connect(cfg)

    e := echo.New()
    e.GET("/", func(c echo.Context) error {
        return c.JSON(200, map[string]string{"message": "FarmFocus API running!"})
    })

	e.GET("tasks", handler.GetTasks)

    e.Logger.Fatal(e.Start(":" + cfg.Port))
}