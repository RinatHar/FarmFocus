package handler

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
	"github.com/labstack/echo/v4"
)

type UserStatHandler struct {
	BaseHandler
	repo *repository.UserStatRepo
}

func NewUserStatHandler(repo *repository.UserStatRepo) *UserStatHandler {
	return &UserStatHandler{repo: repo}
}

// GetUserStats возвращает статистику пользователя
func (h *UserStatHandler) GetUserStats(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	stats, err := h.repo.GetByUserID(context.Background(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			// Создаем статистику если не существует
			stats = &model.UserStat{
				UserID:    userID,
				UpdatedAt: time.Now(),
			}
			if err := h.repo.Create(context.Background(), stats); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
			}
		} else {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, stats)
}

// AddExperience добавляет опыт пользователю
func (h *UserStatHandler) AddExperience(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req struct {
		Amount int64 `json:"amount"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.repo.AddExperience(context.Background(), userID, req.Amount); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Experience added successfully"})
}

// AddGold добавляет золото пользователю
func (h *UserStatHandler) AddGold(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req struct {
		Amount int64 `json:"amount"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.repo.AddGold(context.Background(), userID, req.Amount); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Gold added successfully"})
}

// IncrementStreak увеличивает стрик пользователя
func (h *UserStatHandler) IncrementStreak(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	if err := h.repo.IncrementStreak(context.Background(), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Streak incremented"})
}

// ResetStreak сбрасывает стрик пользователя
func (h *UserStatHandler) ResetStreak(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	if err := h.repo.ResetStreak(context.Background(), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Streak reset"})
}
