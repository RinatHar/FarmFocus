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

type UserPlantHandler struct {
	BaseHandler
	repo         *repository.UserPlantRepo
	userSeedRepo *repository.UserSeedRepo
	userStatRepo *repository.UserStatRepo
	bedRepo      *repository.BedRepo
	seedRepo     *repository.SeedRepo
}

func NewUserPlantHandler(
	repo *repository.UserPlantRepo,
	userSeedRepo *repository.UserSeedRepo,
	userStatRepo *repository.UserStatRepo,
	bedRepo *repository.BedRepo,
	seedRepo *repository.SeedRepo,
) *UserPlantHandler {
	return &UserPlantHandler{
		repo:         repo,
		userSeedRepo: userSeedRepo,
		userStatRepo: userStatRepo,
		bedRepo:      bedRepo,
		seedRepo:     seedRepo,
	}
}

// GetUserPlants godoc
// @Summary Получить все растения пользователя
// @Description Возвращает список всех растений текущего пользователя
// @Tags user-plants
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.UserPlant
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-plants [get]
func (h *UserPlantHandler) GetUserPlants(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	plants, err := h.repo.GetByUser(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, plants)
}

// GetUserPlantByID godoc
// @Summary Получить растение по ID
// @Description Возвращает растение по указанному ID
// @Tags user-plants
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Plant ID"
// @Success 200 {object} model.UserPlant
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-plants/{id} [get]
func (h *UserPlantHandler) GetUserPlantByID(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid plant ID"})
	}

	plant, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Проверяем, что растение принадлежит пользователю
	if plant.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}

	return c.JSON(http.StatusOK, plant)
}

// CreateUserPlant godoc
// @Summary Посадить новое растение
// @Description Сажает новое растение на свободную грядку по cellNumber, используя семя из инвентаря
// @Tags user-plants
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param request body UserPlantCreateRequest true "Данные для посадки растения"
// @Success 200 {object} IPlant
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-plants [post]
func (h *UserPlantHandler) CreateUserPlant(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req UserPlantCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx := context.Background()

	// Ищем грядку по cellNumber у пользователя
	bed, err := h.bedRepo.GetByCellNumber(ctx, userID, req.CellNumber)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Bed not found for this cell number"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Проверяем, что грядка не заблокирована
	if bed.IsLocked {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bed is locked"})
	}

	// Проверяем, что грядка свободна
	existingPlant, err := h.repo.GetByBed(ctx, bed.ID)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if existingPlant != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bed is already occupied"})
	}

	// Проверяем, что у пользователя есть семена
	userSeed, err := h.userSeedRepo.GetByUserAndSeed(ctx, userID, req.SeedID)
	if err != nil && !strings.Contains(err.Error(), "not found") {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if userSeed == nil || userSeed.Quantity < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Not enough seeds"})
	}

	// Используем одно семя
	success, err := h.userSeedRepo.SubtractQuantity(ctx, userID, req.SeedID, 1)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if !success {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to use seed"})
	}

	// Создаем растение
	plant := &model.UserPlant{
		UserID:        userID,
		SeedID:        req.SeedID,
		BedID:         bed.ID, // Используем найденный bed.ID
		CurrentGrowth: 0,
	}

	if err := h.repo.Create(ctx, plant); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем детали семени для формирования IPlant
	seed, err := h.seedRepo.GetByID(ctx, req.SeedID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to get seed details"})
	}

	// Формируем ответ в формате IPlant
	iPlant := IPlant{
		ID:           plant.ID,
		Name:         seed.Name,
		CurrentGrowth: plant.CurrentGrowth,
		TargetGrowth:   seed.TargetGrowth,
		ImgPath:      seed.ImgPlant, // Используем поле img_plant из семени
	}

	return c.JSON(http.StatusOK, iPlant)
}

// AddGrowth godoc
// @Summary Добавить рост растению
// @Description Увеличивает текущий рост растения на указанное количество
// @Tags user-plants
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Plant ID"
// @Param request body UserPlantGrowthRequest true "Количество роста для добавления"
// @Success 200 {object} UserPlantGrowthResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-plants/{id}/add-growth [post]
func (h *UserPlantHandler) AddGrowth(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid plant ID"})
	}

	var req UserPlantGrowthRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Проверяем, что растение принадлежит пользователю
	plant, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if plant.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}

	newGrowth, err := h.repo.AddGrowth(context.Background(), id, req.GrowthAmount)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, UserPlantGrowthResponse{
		NewGrowth: newGrowth,
	})
}

// HarvestPlant godoc
// @Summary Собрать растение
// @Description Собирает готовое к сбору растение и начисляет награды
// @Tags user-plants
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Plant ID"
// @Success 200 {object} model.UserPlantHarvestResult
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-plants/{id}/harvest [post]
func (h *UserPlantHandler) HarvestPlant(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid plant ID"})
	}

	ctx := context.Background()

	// Получаем растение с деталями семени
	plants, err := h.repo.GetWithSeedDetails(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	var plantToHarvest *model.UserPlantWithSeed
	for _, plant := range plants {
		if plant.ID == id {
			plantToHarvest = &plant
			break
		}
	}

	if plantToHarvest == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Plant not found"})
	}

	// Проверяем, готово ли растение к сбору
	if plantToHarvest.CurrentGrowth < plantToHarvest.TargetGrowth {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Plant is not ready for harvest"})
	}

	// Вычисляем награды
	goldEarned := (plantToHarvest.GoldReward * plantToHarvest.GrowthPercent) / 100
	xpEarned := (plantToHarvest.XPReward * plantToHarvest.GrowthPercent) / 100

	// Добавляем награды пользователю
	if goldEarned > 0 {
		if err := h.userStatRepo.AddGold(ctx, userID, int64(goldEarned)); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	if xpEarned > 0 {
		if err := h.userStatRepo.AddExperience(ctx, userID, int64(xpEarned)); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	// Увеличиваем счетчик собранных растений
	stats, err := h.userStatRepo.GetByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	stats.TotalPlantHarvested++
	if err := h.userStatRepo.Update(ctx, stats); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Удаляем растение
	if err := h.repo.Delete(ctx, id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	result := model.UserPlantHarvestResult{
		UserPlantWithSeed: *plantToHarvest,
		GoldEarned:        goldEarned,
		XPEarned:          xpEarned,
		IsReady:           true,
	}

	return c.JSON(http.StatusOK, result)
}

// DeleteUserPlant godoc
// @Summary Удалить растение
// @Description Удаляет растение (например, если оно погибло)
// @Tags user-plants
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Plant ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-plants/{id} [delete]
func (h *UserPlantHandler) DeleteUserPlant(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid plant ID"})
	}

	// Проверяем, что растение принадлежит пользователю
	plant, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if plant.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}

	err = h.repo.Delete(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetPlantsWithDetails godoc
// @Summary Получить растения с деталями
// @Description Возвращает растения с подробной информацией о семенах
// @Tags user-plants
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.UserPlantWithSeed
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-plants/with-details [get]
func (h *UserPlantHandler) GetPlantsWithDetails(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	plants, err := h.repo.GetWithSeedDetails(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, plants)
}

// GetReadyForHarvest godoc
// @Summary Получить растения готовые к сбору
// @Description Возвращает растения, которые достигли максимального роста и готовы к сбору
// @Tags user-plants
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.UserPlantWithSeed
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-plants/ready [get]
func (h *UserPlantHandler) GetReadyForHarvest(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	plants, err := h.repo.GetReadyForHarvest(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, plants)
}

// GetGrowingPlants godoc
// @Summary Получить растущие растения
// @Description Возвращает растения, которые еще не достигли максимального роста
// @Tags user-plants
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.UserPlantWithSeed
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /user-plants/growing [get]
func (h *UserPlantHandler) GetGrowingPlants(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	plants, err := h.repo.GetGrowingPlants(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, plants)
}

// DTO для запросов

// UserPlantCreateRequest представляет запрос на посадку растения
type UserPlantCreateRequest struct {
	SeedID     int `json:"seedId" example:"1"`
	CellNumber int `json:"cellNumber" example:"1"`
}

// UserPlantGrowthRequest представляет запрос на добавление роста
type UserPlantGrowthRequest struct {
	GrowthAmount int `json:"growthAmount" example:"10"`
}

// UserPlantGrowthResponse представляет ответ при добавлении роста
type UserPlantGrowthResponse struct {
	NewGrowth int `json:"newGrowth" example:"50"`
}
