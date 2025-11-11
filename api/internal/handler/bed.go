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

// GetUserBeds возвращает все грядки пользователя
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

// GetBedByID возвращает грядку по ID
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

// GetBedByCellNumber возвращает грядку по номеру ячейки
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

// CreateBed создает новую грядку
func (h *BedHandler) CreateBed(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req struct {
		CellNumber int  `json:"cell_number"`
		IsLocked   bool `json:"is_locked"`
	}

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

// UnlockBed разблокирует грядку
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

// LockBed блокирует грядку
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

// GetAvailableBeds возвращает доступные (разблокированные) грядки
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

// GetEmptyBeds возвращает пустые грядки
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

// GetBedsWithPlants возвращает грядки с информацией о растениях
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

// CreateInitialBeds создает начальные грядки для пользователя
func (h *BedHandler) CreateInitialBeds(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req struct {
		Count int `json:"count"`
	}

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
