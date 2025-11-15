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

// GetUserSeeds godoc
// @Summary Получить семена пользователя
// @Description Возвращает список семян в инвентаре пользователя
// @Tags user-seeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.UserSeed
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-seeds [get]
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

// GetUserSeedsWithDetails godoc
// @Summary Получить семена пользователя с деталями
// @Description Возвращает семена пользователя с подробной информацией о каждом типе семян
// @Tags user-seeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.UserSeedWithDetails
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-seeds/with-details [get]
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

// GetAvailableSeeds godoc
// @Summary Получить доступные семена
// @Description Возвращает все существующие семена с информацией о количестве у пользователя
// @Tags user-seeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param level query int false "Фильтр по уровню семян"
// @Success 200 {array} model.AvailableSeed
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-seeds/available [get]
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

// AddSeed godoc
// @Summary Добавить семена пользователю
// @Description Добавляет новые семена или увеличивает количество существующих семян в инвентаре
// @Tags user-seeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param request body UserSeedAddRequest true "Данные для добавления семян"
// @Success 200 {object} model.UserSeed
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-seeds [post]
func (h *UserSeedHandler) AddSeed(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req UserSeedAddRequest
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

// AddQuantity godoc
// @Summary Добавить количество семян
// @Description Увеличивает количество существующих семян в инвентаре
// @Tags user-seeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param seedId path int true "Seed ID"
// @Param request body UserSeedQuantityRequest true "Количество для добавления"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-seeds/{seedId}/add [post]
func (h *UserSeedHandler) AddQuantity(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	seedID, err := strconv.Atoi(c.Param("seedId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid seed ID"})
	}

	var req UserSeedQuantityRequest
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

// SubtractQuantity godoc
// @Summary Уменьшить количество семян
// @Description Уменьшает количество семян в инвентаре (например, при посадке)
// @Tags user-seeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param seedId path int true "Seed ID"
// @Param request body UserSeedQuantityRequest true "Количество для вычитания"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-seeds/{seedId}/subtract [post]
func (h *UserSeedHandler) SubtractQuantity(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	seedID, err := strconv.Atoi(c.Param("seedId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid seed ID"})
	}

	var req UserSeedQuantityRequest
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

// DeleteUserSeed godoc
// @Summary Удалить семена из инвентаря
// @Description Полностью удаляет семена из инвентаря пользователя
// @Tags user-seeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param seedId path int true "Seed ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-seeds/{seedId} [delete]
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

// GetSeedCount godoc
// @Summary Получить количество уникальных семян
// @Description Возвращает количество различных типов семян в инвентаре пользователя
// @Tags user-seeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} UserSeedCountResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-seeds/count [get]
func (h *UserSeedHandler) GetSeedCount(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	count, err := h.repo.GetTotalSeedCount(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, UserSeedCountResponse{
		TotalSeeds: count,
	})
}

// DTO для запросов

// UserSeedAddRequest представляет запрос на добавление семян
type UserSeedAddRequest struct {
	SeedID   int   `json:"seedId" example:"1"`
	Quantity int64 `json:"quantity" example:"5"`
}

// UserSeedQuantityRequest представляет запрос на изменение количества семян
type UserSeedQuantityRequest struct {
	Amount int64 `json:"amount" example:"3"`
}

// UserSeedCountResponse представляет ответ с количеством семян
type UserSeedCountResponse struct {
	TotalSeeds int `json:"totalSeeds" example:"15"`
}
