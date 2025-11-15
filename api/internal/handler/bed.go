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

type BedHandler struct {
	BaseHandler
	repo *repository.BedRepo
}

func NewBedHandler(repo *repository.BedRepo) *BedHandler {
	return &BedHandler{repo: repo}
}

// GetUserBeds godoc
// @Summary Получить все грядки пользователя
// @Description Возвращает список всех грядок текущего пользователя
// @Tags beds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.Bed
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /beds [get]
func (h *BedHandler) GetUserBeds(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	beds, err := h.repo.GetByUser(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, beds)
}

// GetBedByID godoc
// @Summary Получить грядку по ID
// @Description Возвращает грядку по указанному ID
// @Tags beds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Bed ID"
// @Success 200 {object} model.Bed
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /beds/{id} [get]
func (h *BedHandler) GetBedByID(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid bed ID"})
	}

	bed, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Проверяем, что грядка принадлежит пользователю
	if bed.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}

	return c.JSON(http.StatusOK, bed)
}

// GetBedByCellNumber godoc
// @Summary Получить грядку по номеру ячейки
// @Description Возвращает грядку по номеру ячейки для текущего пользователя
// @Tags beds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param cellNumber path int true "Cell Number"
// @Success 200 {object} model.Bed
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /beds/cell/{cellNumber} [get]
func (h *BedHandler) GetBedByCellNumber(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	cellNumber, err := strconv.Atoi(c.Param("cellNumber"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid cell number"})
	}

	bed, err := h.repo.GetByCellNumber(context.Background(), userID, cellNumber)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, bed)
}

// CreateBed godoc
// @Summary Создать новую грядку
// @Description Создает новую грядку для текущего пользователя
// @Tags beds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param request body BedCreateRequest true "Данные для создания грядки"
// @Success 200 {object} model.Bed
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /beds [post]
func (h *BedHandler) CreateBed(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req BedCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	bed := &model.Bed{
		UserID:     userID,
		CellNumber: req.CellNumber,
		IsLocked:   req.IsLocked,
	}

	if err := h.repo.Create(context.Background(), bed); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, bed)
}

// UnlockBed godoc
// @Summary Разблокировать грядку
// @Description Разблокирует указанную грядку
// @Tags beds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Bed ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /beds/{id}/unlock [post]
func (h *BedHandler) UnlockBed(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid bed ID"})
	}

	// Проверяем, что грядка принадлежит пользователю
	bed, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if bed.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}

	err = h.repo.Unlock(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Bed unlocked successfully"})
}

// LockBed godoc
// @Summary Заблокировать грядку
// @Description Блокирует указанную грядку
// @Tags beds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Bed ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /beds/{id}/lock [post]
func (h *BedHandler) LockBed(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid bed ID"})
	}

	// Проверяем, что грядка принадлежит пользователю
	bed, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if bed.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}

	err = h.repo.Lock(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Bed locked successfully"})
}

// GetAvailableBeds godoc
// @Summary Получить доступные грядки
// @Description Возвращает список разблокированных грядок пользователя
// @Tags beds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.Bed
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /beds/available [get]
func (h *BedHandler) GetAvailableBeds(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	beds, err := h.repo.GetAvailableBeds(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, beds)
}

// GetEmptyBeds godoc
// @Summary Получить пустые грядки
// @Description Возвращает список пустых грядок пользователя
// @Tags beds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.Bed
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /beds/empty [get]
func (h *BedHandler) GetEmptyBeds(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	beds, err := h.repo.GetEmptyBeds(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, beds)
}

// GetBedsWithPlants godoc
// @Summary Получить грядки с растениями
// @Description Возвращает грядки с информацией о посаженных растениях
// @Tags beds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.BedWithUserPlant
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /beds/with-plants [get]
func (h *BedHandler) GetBedsWithPlants(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	beds, err := h.repo.GetWithPlants(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, beds)
}

// CreateInitialBeds godoc
// @Summary Создать начальные грядки
// @Description Создает начальный набор грядок для пользователя
// @Tags beds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param request body BedInitialRequest true "Данные для создания начальных грядок"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /beds/init [post]
func (h *BedHandler) CreateInitialBeds(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req BedInitialRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if req.Count <= 0 {
		req.Count = 5 // значение по умолчанию
	}

	err = h.repo.CreateInitialBeds(context.Background(), userID, req.Count)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Initial beds created successfully"})
}

// DTO для запросов

// BedCreateRequest представляет запрос на создание грядки
type BedCreateRequest struct {
	CellNumber int  `json:"cellNumber" example:"1"`
	IsLocked   bool `json:"isLocked" example:"false"`
}

// BedInitialRequest представляет запрос на создание начальных грядок
type BedInitialRequest struct {
	Count int `json:"count" example:"5"`
}
