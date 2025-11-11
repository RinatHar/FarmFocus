package middleware

import (
	"github.com/labstack/echo/v4"
)

func CORSMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Устанавливаем CORS заголовки
			c.Response().Header().Set("Access-Control-Allow-Origin", "*")
			c.Response().Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, PATCH, OPTIONS")
			c.Response().Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-Requested-With")
			c.Response().Header().Set("Access-Control-Expose-Headers", "Content-Length, Content-Range")
			c.Response().Header().Set("Access-Control-Allow-Credentials", "true")
			c.Response().Header().Set("Access-Control-Max-Age", "86400")

			// Обрабатываем preflight OPTIONS запрос
			if c.Request().Method == "OPTIONS" {
				return c.JSON(200, map[string]string{"status": "ok"})
			}

			return next(c)
		}
	}
}
