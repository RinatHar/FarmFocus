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

// GetAll godoc
// @Summary Получить все семена
// @Description Возвращает полный список всех доступных семян
// @Tags seeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.Seed
// @Failure 500 {object} map[string]string
// @Router /seeds [get]
func (h *SeedHandler) GetAll(c echo.Context) error {
	seeds, err := h.repo.GetAll(context.Background())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, seeds)
}

// GetByID godoc
// @Summary Получить семя по ID
// @Description Возвращает информацию о конкретном семени по его ID
// @Tags seeds
// @Accept json
// @Produce json
// @Param id path int true "Seed ID"
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} model.Seed
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /seeds/{id} [get]
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

// GetByLevel godoc
// @Summary Получить семена по уровню
// @Description Возвращает список семян определенного уровня
// @Tags seeds
// @Accept json
// @Produce json
// @Param level query int true "Уровень семян"
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.Seed
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /seeds/level [get]
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

// GetByRarity godoc
// @Summary Получить семена по редкости
// @Description Возвращает список семян определенной редкости
// @Tags seeds
// @Accept json
// @Produce json
// @Param rarity query string true "Редкость семян (common, rare, epic, legendary)"
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.Seed
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /seeds/rarity [get]
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

// Create godoc
// @Summary Создать новое семя
// @Description Создает новое семя (только для администраторов)
// @Tags seeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param request body model.Seed true "Данные семени"
// @Success 200 {object} model.Seed
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /seeds [post]
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

// Update godoc
// @Summary Обновить семя
// @Description Обновляет информацию о семени (только для администраторов)
// @Tags seeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Seed ID"
// @Param request body model.Seed true "Обновленные данные семени"
// @Success 200 {object} model.Seed
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /seeds/{id} [put]
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

// Delete godoc
// @Summary Удалить семя
// @Description Удаляет семя из системы (только для администраторов)
// @Tags seeds
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Seed ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /seeds/{id} [delete]
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
