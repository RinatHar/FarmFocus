package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/RinatHar/FarmFocus/api/internal/context"
	"github.com/labstack/echo/v4"
)

// AuthMiddleware извлекает userID из заголовков и добавляет в контекст
// Временное решение - позже будет заменено на JWT токены
func AuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Пропускаем Swagger UI, health check и создание пользователя
			path := c.Path()
			if strings.HasPrefix(path, "/swagger") ||
				path == "/health" ||
				(path == "/users" && c.Request().Method == "POST") {
				return next(c)
			}

			// Временное решение: получаем userID из заголовка
			userIDHeader := c.Request().Header.Get("X-User-ID")

			if userIDHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{
					"error": "X-User-ID header is required",
				})
			}

			userID, err := strconv.ParseInt(userIDHeader, 10, 64)
			if err != nil {
				return c.JSON(http.StatusBadRequest, map[string]string{
					"error": "Invalid user ID format",
				})
			}

			// Добавляем userID в контекст
			ctx := context.SetUserID(c.Request().Context(), userID)
			c.SetRequest(c.Request().WithContext(ctx))

			return next(c)
		}
	}
}

// OptionalAuthMiddleware - опциональная аутентификация (для публичных endpoints)
func OptionalAuthMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			userIDHeader := c.Request().Header.Get("X-User-ID")

			if userIDHeader != "" {
				userID, err := strconv.ParseInt(userIDHeader, 10, 64)
				if err == nil {
					ctx := context.SetUserID(c.Request().Context(), userID)
					c.SetRequest(c.Request().WithContext(ctx))
				}
			}

			return next(c)
		}
	}
}
