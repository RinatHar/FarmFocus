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

// GetUserStats godoc
// @Summary Получить статистику пользователя
// @Description Возвращает статистику пользователя. Если статистика не существует - создает новую.
// @Tags user-stats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} model.UserStat
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-stats [get]
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

// GetLevelInfo godoc
// @Summary Получить информацию об уровне
// @Description Возвращает текущий уровень, опыт и прогресс до следующего уровня
// @Tags user-stats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} LevelInfoResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-stats/level-info [get]
func (h *UserStatHandler) GetLevelInfo(c echo.Context) error {
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

	level := stats.Level()
	experienceForNextLevel := stats.ExperienceForNextLevel()

	var progressPercent float64
	if experienceForNextLevel > 0 {
		currentXP := float64(stats.Experience)
		nextLevelXP := float64(experienceForNextLevel)
		progressPercent = (currentXP / nextLevelXP) * 100
	}

	response := LevelInfoResponse{
		Level:                  level,
		Experience:             stats.Experience,
		ExperienceForNextLevel: experienceForNextLevel,
		ProgressPercent:        progressPercent,
	}

	return c.JSON(http.StatusOK, response)
}

// AddExperience godoc
// @Summary Добавить опыт пользователю
// @Description Увеличивает количество опыта пользователя на указанное значение
// @Tags user-stats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param request body UserStatAmountRequest true "Количество опыта для добавления"
// @Success 200 {object} UserStatMessageResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-stats/experience [post]
func (h *UserStatHandler) AddExperience(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req UserStatAmountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "amount must be positive"})
	}

	if err := h.repo.AddExperience(context.Background(), userID, req.Amount); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, UserStatMessageResponse{
		Message: "Experience added successfully",
	})
}

// AddGold godoc
// @Summary Добавить золото пользователю
// @Description Увеличивает количество золота пользователя на указанное значение
// @Tags user-stats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param request body UserStatAmountRequest true "Количество золота для добавления"
// @Success 200 {object} UserStatMessageResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-stats/gold [post]
func (h *UserStatHandler) AddGold(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req UserStatAmountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if req.Amount <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "amount must be positive"})
	}

	if err := h.repo.AddGold(context.Background(), userID, req.Amount); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, UserStatMessageResponse{
		Message: "Gold added successfully",
	})
}

// IncrementStreak godoc
// @Summary Увеличить стрик пользователя
// @Description Увеличивает ежедневный стрик пользователя на 1
// @Tags user-stats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} UserStatMessageResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-stats/streak/increment [post]
func (h *UserStatHandler) IncrementStreak(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	if err := h.repo.IncrementStreak(context.Background(), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, UserStatMessageResponse{
		Message: "Streak incremented",
	})
}

// ResetStreak godoc
// @Summary Сбросить стрик пользователя
// @Description Сбрасывает ежедневный стрик пользователя до 0
// @Tags user-stats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} UserStatMessageResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-stats/streak/reset [post]
func (h *UserStatHandler) ResetStreak(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	if err := h.repo.ResetStreak(context.Background(), userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, UserStatMessageResponse{
		Message: "Streak reset",
	})
}

// DTO для запросов

// UserStatAmountRequest представляет запрос на изменение количества (опыта/золота)
type UserStatAmountRequest struct {
	Amount int64 `json:"amount" example:"100"`
}

// UserStatMessageResponse представляет ответ с сообщением
type UserStatMessageResponse struct {
	Message string `json:"message" example:"Operation completed successfully"`
}

// LevelInfoResponse представляет информацию об уровне
type LevelInfoResponse struct {
	Level                  int     `json:"level"`
	Experience             int64   `json:"experience"`
	ExperienceForNextLevel int64   `json:"experienceForNextLevel"`
	ProgressPercent        float64 `json:"progressPercent"`
}
