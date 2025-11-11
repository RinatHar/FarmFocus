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

type UserHandler struct {
	BaseHandler
	repo     *repository.UserRepo
	statRepo *repository.UserStatRepo
	bedRepo  *repository.BedRepo
}

func NewUserHandler(repo *repository.UserRepo, statRepo *repository.UserStatRepo, bedRepo *repository.BedRepo) *UserHandler {
	return &UserHandler{
		repo:     repo,
		statRepo: statRepo,
		bedRepo:  bedRepo,
	}
}

// GetCurrentUser возвращает текущего пользователя
func (h *UserHandler) GetCurrentUser(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	user, err := h.repo.GetByID(context.Background(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}

// CreateOrUpdateUser создает или обновляет пользователя
func (h *UserHandler) CreateOrUpdateUser(c echo.Context) error {
	var req struct {
		MaxID    int64  `json:"max_id"`
		Username string `json:"username"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx := context.Background()

	// Проверяем существует ли пользователь
	existingUser, err := h.repo.GetByMaxID(ctx, req.MaxID)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if existingUser != nil {
		// Обновляем существующего пользователя
		existingUser.Username = req.Username
		now := time.Now()
		existingUser.LastLogin = &now

		if err := h.repo.Update(ctx, existingUser); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		return c.JSON(http.StatusOK, existingUser)
	}

	// Создаем нового пользователя
	user := &model.User{
		MaxID:     req.MaxID,
		Username:  req.Username,
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	if err := h.repo.Create(ctx, user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Создаем статистику пользователя
	stat := &model.UserStat{
		UserID:     user.ID,
		Experience: 0,
		Gold:       100, // начальный капитал
		Streak:     0,
		UpdatedAt:  time.Now(),
	}

	if err := h.statRepo.Create(ctx, stat); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Создаем начальные грядки
	if err := h.bedRepo.CreateInitialBeds(ctx, user.ID, 9); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

// UpdateUser обновляет данные пользователя
func (h *UserHandler) UpdateUser(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req struct {
		Username string `json:"username"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	user, err := h.repo.GetByID(context.Background(), userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	user.Username = req.Username
	now := time.Now()
	user.LastLogin = &now

	if err := h.repo.Update(context.Background(), user); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, user)
}
