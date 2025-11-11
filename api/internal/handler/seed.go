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

type SeedHandler struct {
	BaseHandler
	repo *repository.SeedRepo
}

func NewSeedHandler(repo *repository.SeedRepo) *SeedHandler {
	return &SeedHandler{repo: repo}
}

// GetAll возвращает все семена
func (h *SeedHandler) GetAll(c echo.Context) error {
	seeds, err := h.repo.GetAll(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, seeds)
}

// GetByID возвращает семя по ID
func (h *SeedHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid seed ID"})
	}

	seed, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, seed)
}

// GetByLevel возвращает семена по уровню
func (h *SeedHandler) GetByLevel(c echo.Context) error {
	level, err := strconv.Atoi(c.QueryParam("level"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid level parameter"})
	}

	seeds, err := h.repo.GetByLevel(context.Background(), level)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, seeds)
}

// GetByRarity возвращает семена по редкости
func (h *SeedHandler) GetByRarity(c echo.Context) error {
	rarity := c.QueryParam("rarity")
	if rarity == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "rarity parameter is required"})
	}

	seeds, err := h.repo.GetByRarity(context.Background(), rarity)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, seeds)
}

// Create создает новое семя (админ)
func (h *SeedHandler) Create(c echo.Context) error {
	var seed model.Seed
	if err := c.Bind(&seed); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if err := h.repo.Create(context.Background(), &seed); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, seed)
}

// Update обновляет семя (админ)
func (h *SeedHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid seed ID"})
	}

	var seed model.Seed
	if err := c.Bind(&seed); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	seed.ID = id

	if err := h.repo.Update(context.Background(), &seed); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, seed)
}

// Delete удаляет семя (админ)
func (h *SeedHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid seed ID"})
	}

	err = h.repo.Delete(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
