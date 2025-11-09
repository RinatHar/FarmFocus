package handler

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func GetTasks(c echo.Context) error {
	// Заглушка, позже подключим БД
	tasks := []map[string]interface{}{
		{"id": 1, "title": "First Task"},
		{"id": 2, "title": "Second Task"},
	}
	return c.JSON(http.StatusOK, tasks)
}
