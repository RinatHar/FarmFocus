package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
	"github.com/labstack/echo/v4"
)

type UserSeedHandler struct {
	BaseHandler
	repo *repository.UserSeedRepo
}

func NewUserSeedHandler(repo *repository.UserSeedRepo) *UserSeedHandler {
	return &UserSeedHandler{repo: repo}
}

// GetUserSeeds возвращает семена пользователя
func (h *UserSeedHandler) GetUserSeeds(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	userSeeds, err := h.repo.GetByUser(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, userSeeds)
}

// GetUserSeedsWithDetails возвращает семена пользователя с деталями
func (h *UserSeedHandler) GetUserSeedsWithDetails(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	userSeeds, err := h.repo.GetUserSeedsWithDetails(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, userSeeds)
}

// GetAvailableSeeds возвращает все семена с информацией о количестве у пользователя
func (h *UserSeedHandler) GetAvailableSeeds(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	level, _ := strconv.Atoi(c.QueryParam("level")) // уровень опциональный

	seeds, err := h.repo.GetAvailableSeedsForUser(context.Background(), userID, level)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, seeds)
}

// AddSeed добавляет семена пользователю
func (h *UserSeedHandler) AddSeed(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req struct {
		SeedID   int   `json:"seed_id"`
		Quantity int64 `json:"quantity"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	userSeed := &model.UserSeed{
		UserID:   userID,
		SeedID:   req.SeedID,
		Quantity: req.Quantity,
	}

	if err := h.repo.CreateOrUpdate(context.Background(), userSeed); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, userSeed)
}

// AddQuantity добавляет количество к существующим семенам
func (h *UserSeedHandler) AddQuantity(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	seedID, err := strconv.Atoi(c.Param("seedId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid seed ID"})
	}

	var req struct {
		Amount int64 `json:"amount"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	err = h.repo.AddQuantity(context.Background(), userID, seedID, req.Amount)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// SubtractQuantity вычитает количество семян
func (h *UserSeedHandler) SubtractQuantity(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	seedID, err := strconv.Atoi(c.Param("seedId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid seed ID"})
	}

	var req struct {
		Amount int64 `json:"amount"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	success, err := h.repo.SubtractQuantity(context.Background(), userID, seedID, req.Amount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if !success {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "insufficient seed quantity"})
	}

	return c.NoContent(http.StatusNoContent)
}

// DeleteUserSeed удаляет запись о семенах пользователя
func (h *UserSeedHandler) DeleteUserSeed(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	seedID, err := strconv.Atoi(c.Param("seedId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid seed ID"})
	}

	err = h.repo.Delete(context.Background(), userID, seedID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetSeedCount возвращает количество уникальных семян у пользователя
func (h *UserSeedHandler) GetSeedCount(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	count, err := h.repo.GetTotalSeedCount(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]int{"total_seeds": count})
}
