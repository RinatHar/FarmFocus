package handler

import (
	"github.com/RinatHar/FarmFocus/api/internal/context"
	"github.com/labstack/echo/v4"
)

// BaseHandler содержит общие методы для всех хендлеров
type BaseHandler struct{}

// GetUserIDFromContext извлекает userID из контекста Echo
func (h *BaseHandler) GetUserIDFromContext(c echo.Context) (int64, error) {
	userID, ok := context.GetUserID(c.Request().Context())
	if !ok {
		return 0, echo.NewHTTPError(401, "User not authenticated")
	}
	return userID, nil
}

// MustGetUserIDFromContext извлекает userID из контекста (паникует если не найден)
func (h *BaseHandler) MustGetUserIDFromContext(c echo.Context) int64 {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		panic(err)
	}
	return userID
}
