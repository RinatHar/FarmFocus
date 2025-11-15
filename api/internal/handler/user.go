package handler

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
	"github.com/labstack/echo/v4"
)

type UserHandler struct {
	BaseHandler
	repo          *repository.UserRepo
	statRepo      *repository.UserStatRepo
	bedRepo       *repository.BedRepo
	taskRepo      *repository.TaskRepo
	habitRepo     *repository.HabitRepo
	tagRepo       *repository.TagRepo
	seedRepo      *repository.SeedRepo
	userSeedRepo  *repository.UserSeedRepo
	userPlantRepo *repository.UserPlantRepo
	goodRepo      *repository.GoodRepo
	progressRepo  *repository.ProgressLogRepo
}

func NewUserHandler(
	repo *repository.UserRepo,
	statRepo *repository.UserStatRepo,
	bedRepo *repository.BedRepo,
	taskRepo *repository.TaskRepo,
	habitRepo *repository.HabitRepo,
	tagRepo *repository.TagRepo,
	seedRepo *repository.SeedRepo,
	userSeedRepo *repository.UserSeedRepo,
	userPlantRepo *repository.UserPlantRepo,
	goodRepo *repository.GoodRepo,
	progressRepo *repository.ProgressLogRepo,
) *UserHandler {
	return &UserHandler{
		repo:          repo,
		statRepo:      statRepo,
		bedRepo:       bedRepo,
		taskRepo:      taskRepo,
		habitRepo:     habitRepo,
		tagRepo:       tagRepo,
		seedRepo:      seedRepo,
		userSeedRepo:  userSeedRepo,
		userPlantRepo: userPlantRepo,
		goodRepo:      goodRepo,
		progressRepo:  progressRepo,
	}
}

// RecoverPlants godoc
// @Summary Восстановить засохшие растения
// @Description Восстанавливает засохшие растения после выполнения ежедневного задания
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/recover-plants [post]
func (h *UserHandler) RecoverPlants(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Проверяем, выполнено ли сегодня задание
	hasCompletedToday, err := h.progressRepo.HasUserCompletedTaskToday(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	if !hasCompletedToday {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Необходимо выполнить хотя бы одно задание сегодня для восстановления растений",
		})
	}

	// Восстанавливаем растения
	if err := h.userPlantRepo.ResetWitheredStatus(ctx, userID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Растения успешно восстановлены",
	})
}

// SyncUserData godoc
// @Summary Синхронизировать все данные пользователя
// @Description Возвращает все данные пользователя для синхронизации в формате ServerFarmData. Если пользователь не существует - создает его с дефолтными данными.
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} ServerFarmData
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/sync [get]
func (h *UserHandler) SyncUserData(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Проверяем существует ли пользователь, если нет - создаем через CreateOrUpdateUser
	_, err = h.repo.GetByID(ctx, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			// Создаем пользователя с дефолтными данными
			if err := h.createUserWithDefaults(ctx, userID); err != nil {
				return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create user: " + err.Error()})
			}
		}
	}

	// Остальная логика синхронизации без изменений...
	stats, err := h.statRepo.GetByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// ... остальной код синхронизации без изменений
	didTaskToday, err := h.progressRepo.HasUserCompletedTaskToday(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	isDrought := stats.IsDrought
	tasks, err := h.taskRepo.GetAll(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	habits, err := h.habitRepo.GetAll(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	tags, err := h.tagRepo.GetByUser(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	beds, err := h.bedRepo.GetByUser(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	inventorySeeds, err := h.userSeedRepo.GetUserSeedsWithDetails(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	shopGoods, err := h.goodRepo.GetByUser(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	shopStorage, err := h.transformShopGoods(ctx, shopGoods)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	plants, err := h.userPlantRepo.GetWithSeedDetails(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	plantMap := make(map[int]IPlant)
	for _, plant := range plants {
		plantMap[plant.BedID] = IPlant{
			ID:            plant.ID,
			Name:          plant.SeedName,
			CurrentGrowth: plant.CurrentGrowth,
			TargetGrowth:  plant.TargetGrowth,
			ImgPath:       plant.SeedImgPlant,
		}
	}

	field := make([]IBed, len(beds))
	for i, bed := range beds {
		var plant *IPlant
		if p, exists := plantMap[bed.ID]; exists {
			plant = &p
		}

		field[i] = IBed{
			ID:     bed.ID,
			Plant:  plant,
			IsLock: bed.IsLocked,
		}
	}

	seeds := make([]SeedStorage, len(plants))
	for i, plant := range plants {
		seeds[i] = SeedStorage{
			ID:            plant.ID,
			SeedID:        plant.SeedID,
			CurrentGrowth: plant.CurrentGrowth,
			IsWithered:    plant.IsWithered,
			BedID:         plant.BedID,
			CreatedAt:     plant.CreatedAt,
			SeedName:      plant.SeedName,
			TargetGrowth:  plant.TargetGrowth,
			Icon:          plant.SeedIcon,
			ImgPath:       plant.SeedImgPlant,
		}
	}

	response := ServerFarmData{
		CurrentXp:      int(stats.Experience),
		Coins:          int(stats.Gold),
		Strick:         stats.CurrentStreak,
		DidTaskToday:   didTaskToday,
		IsDrought:      isDrought,
		Tasks:          tasks,
		Habits:         habits,
		Tags:           tags,
		Field:          field,
		Seeds:          seeds,
		InventorySeeds: inventorySeeds,
		ShopStorage:    shopStorage,
	}

	return c.JSON(http.StatusOK, response)
}

// createUserWithDefaults создает пользователя с дефолтными данными
func (h *UserHandler) createUserWithDefaults(ctx context.Context, userID int64) error {
	// Создаем пользователя
	user := &model.User{
		MaxID:     userID,
		Username:  fmt.Sprintf("User_%d", userID),
		CreatedAt: time.Now(),
		IsActive:  true,
	}

	if err := h.repo.Create(ctx, user); err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	// Создаем статистику пользователя
	stat := &model.UserStat{
		UserID:              user.ID,
		Experience:          0,
		Gold:                10,
		CurrentStreak:       0,
		LongestStreak:       0,
		TotalTasksCompleted: 0,
		TotalPlantHarvested: 0,
		UpdatedAt:           time.Now(),
	}

	if err := h.statRepo.Create(ctx, stat); err != nil {
		return fmt.Errorf("failed to create user stat: %w", err)
	}

	// Создаем начальные грядки
	if err := h.bedRepo.CreateInitialBeds(ctx, user.ID, 9); err != nil {
		return fmt.Errorf("failed to create initial beds: %w", err)
	}

	// Добавляем стартовый набор семян (10 пшеницы)
	wheatSeedID := 1
	initialWheatQuantity := 10
	if err := h.userSeedRepo.AddOrUpdateQuantity(ctx, user.ID, wheatSeedID, initialWheatQuantity); err != nil {
		return fmt.Errorf("failed to add initial seeds: %w", err)
	}

	// Создаем стартовый набор товаров в магазине
	if err := h.createInitialShopGoods(ctx, user.ID); err != nil {
		return fmt.Errorf("failed to create initial shop goods: %w", err)
	}

	return nil
}

// transformShopGoods преобразует товары из магазина в нужный формат
func (h *UserHandler) transformShopGoods(ctx context.Context, goods []model.Good) ([]ShopItem, error) {
	var shopItems []ShopItem

	for _, good := range goods {
		shopItem := ShopItem{
			ID:       good.ID,
			Quantity: good.Quantity,
			Cost:     good.Cost,
			Type:     good.Type,
		}

		switch good.Type {
		case "seed":
			// Для семян получаем детали семени
			seed, err := h.seedRepo.GetByID(ctx, good.IDGood)
			if err != nil {
				return nil, fmt.Errorf("failed to get seed details for ID %d: %w", good.IDGood, err)
			}
			shopItem.Item = SeedItem{
				ID:           seed.ID,
				Name:         seed.Name,
				Icon:         seed.Icon,
				ImgPlant:     seed.ImgPlant,
				TargetGrowth: seed.TargetGrowth,
				Rarity:       seed.Rarity,
			}

		default:
			// Для других типов товаров можно добавить соответствующую логику
			shopItem.Item = nil
		}

		shopItems = append(shopItems, shopItem)
	}

	return shopItems, nil
}

// shouldHabitBeCompletedYesterday проверяет, должна ли привычка быть выполнена вчера
func (h *UserHandler) shouldHabitBeCompletedYesterday(habit model.Habit, yesterday time.Time) bool {
	// Для ежедневных привычек - должны выполняться каждый день
	if habit.Period == "day" {
		return true
	}

	// Для еженедельных - если вчера был день выполнения
	if habit.Period == "week" {
		daysSinceStart := int(yesterday.Sub(habit.StartDate).Hours() / 24)
		return daysSinceStart%7 == 0
	}

	// Для ежемесячных - если вчера был день месяца
	if habit.Period == "month" {
		return yesterday.Day() == habit.StartDate.Day()
	}

	return false
}

// GetCurrentUser godoc
// @Summary Получить текущего пользователя
// @Description Возвращает информацию о текущем аутентифицированном пользователе
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} model.User
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [get]
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

// CreateOrUpdateUser godoc
// @Summary Создать или обновить пользователя
// @Description Создает нового пользователя или обновляет существующего по MaxID. При создании также создается статистика, начальные грядки, добавляется стартовый набор семян и товары в магазин.
// @Tags users
// @Accept json
// @Produce json
// @Param request body UserCreateRequest true "Данные пользователя"
// @Success 200 {object} model.User "Пользователь обновлен"
// @Success 201 {object} model.User "Пользователь создан"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *UserHandler) CreateOrUpdateUser(c echo.Context) error {
	var req UserCreateRequest
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
		UserID:              user.ID,
		Experience:          0,
		Gold:                10,
		CurrentStreak:       0,
		LongestStreak:       0,
		TotalTasksCompleted: 0,
		TotalPlantHarvested: 0,
		UpdatedAt:           time.Now(),
	}

	if err := h.statRepo.Create(ctx, stat); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Создаем начальные грядки
	if err := h.bedRepo.CreateInitialBeds(ctx, user.ID, 9); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Добавляем стартовый набор семян (10 пшеницы)
	wheatSeedID := 1 // ID пшеницы из таблицы seed
	initialWheatQuantity := 10

	if err := h.userSeedRepo.AddOrUpdateQuantity(ctx, user.ID, wheatSeedID, initialWheatQuantity); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to add initial seeds: " + err.Error()})
	}

	// Добавляем стартовый набор товаров в магазин
	if err := h.createInitialShopGoods(ctx, user.ID); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create initial shop goods: " + err.Error()})
	}

	return c.JSON(http.StatusCreated, user)
}

// createInitialShopGoods создает начальный набор товаров в магазине пользователя
func (h *UserHandler) createInitialShopGoods(ctx context.Context, userID int64) error {
	aubergineGood := &model.Good{
		UserID:   userID,
		Type:     "seed",
		IDGood:   2,
		Quantity: 5,
		Cost:     4,
	}

	if err := h.goodRepo.CreateOrUpdate(ctx, aubergineGood); err != nil {
		return fmt.Errorf("failed to create aubergine good: %w", err)
	}

	wheatGood := &model.Good{
		UserID:   userID,
		Type:     "seed",
		IDGood:   1,
		Quantity: 10,
		Cost:     1,
	}

	if err := h.goodRepo.CreateOrUpdate(ctx, wheatGood); err != nil {
		return fmt.Errorf("failed to create wheat good: %w", err)
	}

	// 1 дополнительное поле (bed) по цене 100 золота
	bedGood := &model.Good{
		UserID:   userID,
		Type:     "bed",
		IDGood:   0,
		Quantity: 1,
		Cost:     30,
	}

	if err := h.goodRepo.CreateOrUpdate(ctx, bedGood); err != nil {
		return fmt.Errorf("failed to create bed good: %w", err)
	}

	return nil
}

// UpdateUser godoc
// @Summary Обновить данные пользователя
// @Description Обновляет информацию о текущем пользователе (например, имя пользователя)
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param request body UserUpdateRequest true "Обновленные данные пользователя"
// @Success 200 {object} model.User
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [put]
func (h *UserHandler) UpdateUser(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req UserUpdateRequest
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

// DTO для запросов

// UserCreateRequest представляет запрос на создание/обновление пользователя
type UserCreateRequest struct {
	MaxID    int64  `json:"maxId" example:"123456789"`
	Username string `json:"username" example:"john_doe"`
}

// UserUpdateRequest представляет запрос на обновление пользователя
type UserUpdateRequest struct {
	Username string `json:"username" example:"john_doe_updated"`
}

// IPlant представляет растение на грядке
type IPlant struct {
	ID            int    `json:"id"`
	Name          string `json:"name"`
	CurrentGrowth int    `json:"currentGrowth"`
	TargetGrowth  int    `json:"targetGrowth"`
	ImgPath       string `json:"imgPath"`
}

// IBed представляет грядку с растением
type IBed struct {
	ID     int     `json:"id"`
	Plant  *IPlant `json:"plant"`
	IsLock bool    `json:"isLock"`
}

// ServerFarmData основной формат ответа для синхронизации
type ServerFarmData struct {
	CurrentXp      int                         `json:"currentXp"`
	Coins          int                         `json:"coins"`
	Strick         int                         `json:"strick"`
	DidTaskToday   bool                        `json:"didTaskToday"`
	IsDrought      bool                        `json:"isDrought"`
	Tasks          []model.Task                `json:"tasks"`
	Habits         []model.Habit               `json:"habits"`
	Tags           []model.Tag                 `json:"tags"`
	Field          []IBed                      `json:"field"`
	Seeds          []SeedStorage               `json:"seeds"`
	InventorySeeds []model.UserSeedWithDetails `json:"inventorySeeds"`
	ShopStorage    []ShopItem                  `json:"shopItem"`
}

// ShopItem представляет товар в магазине
type ShopItem struct {
	ID       int         `json:"id"`
	Quantity int         `json:"quantity"`
	Cost     int         `json:"cost"`
	Type     string      `json:"type"`
	Item     interface{} `json:"item,omitempty"`
}

// SeedItem представляет семя в магазине
type SeedItem struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Icon         string `json:"icon"`
	ImgPlant     string `json:"imgPlant"`
	TargetGrowth int    `json:"targetGrowth"`
	Rarity       string `json:"rarity"`
}

// SeedStorage представляет растение на грядке (для отдельного списка)
type SeedStorage struct {
	ID            int       `json:"id"`
	SeedID        int       `json:"seedId"`
	CurrentGrowth int       `json:"currentGrowth"`
	IsWithered    bool      `json:"isWithered"`
	BedID         int       `json:"bedId"`
	CreatedAt     time.Time `json:"createdAt"`
	SeedName      string    `json:"seedName,omitempty"`
	TargetGrowth  int       `json:"targetGrowth,omitempty"`
	Icon          string    `json:"icon,omitempty"`
	ImgPath       string    `json:"imgPath,omitempty"`
}
