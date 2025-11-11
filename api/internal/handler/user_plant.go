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

// GetUserPlants возвращает все растения пользователя
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

// GetUserPlantByID возвращает растение по ID
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

// CreateUserPlant создает новое растение
func (h *UserPlantHandler) CreateUserPlant(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req struct {
		SeedID int `json:"seed_id"`
		BedID  int `json:"bed_id"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	ctx := context.Background()

	// Проверяем, что грядка принадлежит пользователю и свободна
	bed, err := h.bedRepo.GetByID(ctx, req.BedID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "Bed not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if bed.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "Bed does not belong to user"})
	}

	if bed.IsLocked {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Bed is locked"})
	}

	// Проверяем, что грядка свободна
	existingPlant, err := h.repo.GetByBed(ctx, req.BedID)
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
		BedID:         req.BedID,
		CurrentGrowth: 0,
	}

	if err := h.repo.Create(ctx, plant); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, plant)
}

// AddGrowth добавляет рост растению
func (h *UserPlantHandler) AddGrowth(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid plant ID"})
	}

	var req model.UserPlantGrowthRequest
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

	return c.JSON(http.StatusOK, map[string]int{"new_growth": newGrowth})
}

// HarvestPlant собирает растение
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

// DeleteUserPlant удаляет растение
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

// GetPlantsWithDetails возвращает растения с деталями семян
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

// GetReadyForHarvest возвращает растения готовые к сбору
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

// GetGrowingPlants возвращает растущие растения
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
