package handler

import (
	"context"
	"net/http"
	"strconv"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
	"github.com/labstack/echo/v4"
)

type GoodHandler struct {
	BaseHandler
	repo         *repository.GoodRepo
	userStatRepo *repository.UserStatRepo
	userSeedRepo *repository.UserSeedRepo
	bedRepo      *repository.BedRepo
	seedRepo     *repository.SeedRepo
}

func NewGoodHandler(
	repo *repository.GoodRepo,
	userStatRepo *repository.UserStatRepo,
	userSeedRepo *repository.UserSeedRepo,
	bedRepo *repository.BedRepo,
	seedRepo *repository.SeedRepo,
) *GoodHandler {
	return &GoodHandler{
		repo:         repo,
		userStatRepo: userStatRepo,
		userSeedRepo: userSeedRepo,
		bedRepo:      bedRepo,
		seedRepo:     seedRepo,
	}
}

// GetUserGoods godoc
// @Summary Получить товары пользователя
// @Description Возвращает список всех товаров текущего пользователя
// @Tags goods
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.Good
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /goods [get]
func (h *GoodHandler) GetUserGoods(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	goods, err := h.repo.GetByUser(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, goods)
}

// GetUserGoodsByType godoc
// @Summary Получить товары пользователя по типу
// @Description Возвращает товары указанного типа для текущего пользователя
// @Tags goods
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param type path string true "Good Type (seed/bed/tool/fertilizer)"
// @Success 200 {array} model.Good
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /goods/type/{type} [get]
func (h *GoodHandler) GetUserGoodsByType(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	goodType := c.Param("type")
	if !isValidGoodType(goodType) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid good type. Must be: seed, bed, tool, or fertilizer",
		})
	}

	goods, err := h.repo.GetByUserAndType(context.Background(), userID, goodType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, goods)
}

// CreateGood godoc
// @Summary Создать товар
// @Description Создает новый товар для пользователя
// @Tags goods
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param request body CreateGoodRequest true "Данные товара"
// @Success 201 {object} model.Good
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /goods [post]
func (h *GoodHandler) CreateGood(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req CreateGoodRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	if !isValidGoodType(req.Type) {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid good type. Must be: seed, bed, tool, or fertilizer",
		})
	}

	good := model.Good{
		UserID:   userID,
		Type:     req.Type,
		IDGood:   req.IDGood,
		Quantity: req.Quantity,
		Cost:     req.Cost,
	}

	if err := h.repo.Create(context.Background(), &good); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, good)
}

// UpdateGoodQuantity godoc
// @Summary Обновить количество товара
// @Description Обновляет количество товара по ID
// @Tags goods
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Good ID"
// @Param request body UpdateQuantityRequest true "Новое количество"
// @Success 200 {object} model.Good
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /goods/{id}/quantity [put]
func (h *GoodHandler) UpdateGoodQuantity(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid good ID"})
	}

	var req UpdateQuantityRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Проверяем, что товар принадлежит пользователю
	good, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	if good.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}

	if err := h.repo.UpdateQuantity(context.Background(), id, req.Quantity); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	updatedGood, _ := h.repo.GetByID(context.Background(), id)
	return c.JSON(http.StatusOK, updatedGood)
}

// UpdateGoodCost godoc
// @Summary Обновить цену товара
// @Description Обновляет цену товара по ID
// @Tags goods
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Good ID"
// @Param request body UpdateCostRequest true "Новая цена"
// @Success 200 {object} model.Good
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /goods/{id}/cost [put]
func (h *GoodHandler) UpdateGoodCost(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid good ID"})
	}

	var req UpdateCostRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Проверяем, что товар принадлежит пользователю
	good, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	if good.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}

	if err := h.repo.UpdateCost(context.Background(), id, req.Cost); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	updatedGood, _ := h.repo.GetByID(context.Background(), id)
	return c.JSON(http.StatusOK, updatedGood)
}

// AddQuantity godoc
// @Summary Добавить количество товара
// @Description Добавляет указанное количество к существующему товару
// @Tags goods
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Good ID"
// @Param request body AddQuantityRequest true "Количество для добавления"
// @Success 200 {object} model.Good
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /goods/{id}/add-quantity [patch]
func (h *GoodHandler) AddQuantity(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid good ID"})
	}

	var req AddQuantityRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Проверяем, что товар принадлежит пользователю
	good, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	if good.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}

	if err := h.repo.AddQuantity(context.Background(), id, req.Amount); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	updatedGood, _ := h.repo.GetByID(context.Background(), id)
	return c.JSON(http.StatusOK, updatedGood)
}

// DeleteGood godoc
// @Summary Удалить товар
// @Description Удаляет товар по ID
// @Tags goods
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Good ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /goods/{id} [delete]
func (h *GoodHandler) DeleteGood(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid good ID"})
	}

	// Проверяем, что товар принадлежит пользователю
	good, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	if good.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}

	if err := h.repo.Delete(context.Background(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// CreateBatchGoods godoc
// @Summary Создать несколько товаров
// @Description Создает несколько товаров для пользователя
// @Tags goods
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param request body []CreateGoodRequest true "Массив товаров"
// @Success 201 {array} model.Good
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /goods/batch [post]
func (h *GoodHandler) CreateBatchGoods(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req []CreateGoodRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	var createdGoods []model.Good
	for _, goodReq := range req {
		if !isValidGoodType(goodReq.Type) {
			continue // Пропускаем невалидные типы
		}

		good := model.Good{
			UserID:   userID,
			Type:     goodReq.Type,
			IDGood:   goodReq.IDGood,
			Quantity: goodReq.Quantity,
			Cost:     goodReq.Cost,
		}

		if err := h.repo.CreateOrUpdate(context.Background(), &good); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}

		createdGoods = append(createdGoods, good)
	}

	return c.JSON(http.StatusCreated, createdGoods)
}

// BuyGood godoc
// @Summary Купить товар
// @Description Покупает один товар из магазина, списывая золото и добавляя товар в инвентарь
// @Tags goods
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Good ID"
// @Success 200 {object} BuyGoodResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /goods/{id}/buy [post]
func (h *GoodHandler) BuyGood(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	// Используем "id" как в роутере
	goodID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid good ID"})
	}

	ctx := context.Background()

	// Получаем товар
	good, err := h.repo.GetByID(ctx, goodID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "good not found"})
	}

	// Проверяем, что товар принадлежит пользователю (магазин пользователя)
	if good.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}

	// Проверяем доступное количество
	if good.Quantity < 1 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "good is out of stock",
		})
	}

	// Получаем статистику пользователя для проверки золота
	stat, err := h.userStatRepo.GetByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Проверяем достаточно ли золота
	if stat.Gold < int64(good.Cost) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "not enough gold"})
	}

	// Списываем золото
	if err := h.userStatRepo.AddGold(ctx, userID, -int64(good.Cost)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to deduct gold"})
	}

	// Уменьшаем количество товара в магазине на 1
	newQuantity := good.Quantity - 1
	if err := h.repo.UpdateQuantity(ctx, goodID, newQuantity); err != nil {
		// Если не удалось обновить количество, возвращаем золото
		h.userStatRepo.AddGold(ctx, userID, int64(good.Cost))
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to update good quantity"})
	}

	var purchasedItem interface{}

	// Добавляем товар в инвентарь пользователя в зависимости от типа
	switch good.Type {
	case "seed":
		// Добавляем семена в инвентарь
		if err := h.userSeedRepo.AddOrUpdateQuantity(ctx, userID, good.IDGood, 1); err != nil {
			// Если не удалось добавить в инвентарь, возвращаем товар и золото
			h.repo.UpdateQuantity(ctx, goodID, good.Quantity)
			h.userStatRepo.AddGold(ctx, userID, int64(good.Cost))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to add seeds to inventory"})
		}

		// Получаем детали семени для ответа
		seed, err := h.seedRepo.GetByID(ctx, good.IDGood)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to get seed details"})
		}

		purchasedItem = SeedItem{
			ID:           seed.ID,
			Name:         seed.Name,
			Icon:         seed.Icon,
			ImgPlant:     seed.ImgPlant,
			TargetGrowth: seed.TargetGrowth,
			Rarity:       seed.Rarity,
		}

	case "bed":
		// Разблокируем грядку
		unlockedBed, err := h.bedRepo.UnlockNextBed(ctx, userID)
		if err != nil {
			// Если не удалось разблокировать, возвращаем товар и золото
			h.repo.UpdateQuantity(ctx, goodID, good.Quantity)
			h.userStatRepo.AddGold(ctx, userID, int64(good.Cost))
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to unlock bed"})
		}

		purchasedItem = UnlockedBedInfo{
			ID:         unlockedBed.ID,
			CellNumber: unlockedBed.CellNumber,
			IsLocked:   unlockedBed.IsLocked,
		}

	default:
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "unsupported good type"})
	}

	// Получаем обновленный товар
	updatedGood, _ := h.repo.GetByID(ctx, goodID)

	return c.JSON(http.StatusOK, BuyGoodResponse{
		GoodID:        goodID,
		Quantity:      1,
		TotalCost:     good.Cost,
		Remaining:     updatedGood.Quantity,
		Message:       "Purchase successful",
		ItemType:      good.Type,
		ItemID:        good.IDGood,
		PurchasedItem: purchasedItem,
	})
}

// Вспомогательные функции и DTO

func isValidGoodType(goodType string) bool {
	validTypes := map[string]bool{
		"seed": true, "bed": true, "tool": true, "fertilizer": true,
	}
	return validTypes[goodType]
}

type BuyGoodResponse struct {
	GoodID        int         `json:"goodId"`
	Quantity      int         `json:"quantity"`
	TotalCost     int         `json:"totalCost"`
	Remaining     int         `json:"remaining"`
	Message       string      `json:"message"`
	ItemType      string      `json:"itemType"`
	ItemID        int         `json:"itemId"`
	PurchasedItem interface{} `json:"purchasedItem"` // Добавлено поле с купленным предметом
}

type UnlockedBedInfo struct {
	ID         int  `json:"id"`
	CellNumber int  `json:"cellNumber"`
	IsLocked   bool `json:"isLocked"`
}

type CreateGoodRequest struct {
	Type     string `json:"type" binding:"required"`
	IDGood   int    `json:"idGood" binding:"required,min=1"`
	Quantity int    `json:"quantity" binding:"min=0"`
	Cost     int    `json:"cost" binding:"required,min=0"`
}

type UpdateQuantityRequest struct {
	Quantity int `json:"quantity" binding:"required,min=0"`
}

type UpdateCostRequest struct {
	Cost int `json:"cost" binding:"required,min=0"`
}

type AddQuantityRequest struct {
	Amount int `json:"amount" binding:"required,min=1"`
}
