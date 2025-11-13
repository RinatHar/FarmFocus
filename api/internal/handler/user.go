package handler

import (
	"context"
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
// @Description Возвращает все данные пользователя для синхронизации в формате ServerFarmData
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

	// Получаем статистику пользователя
	stats, err := h.statRepo.GetByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Проверяем, выполнял ли пользователь задачи сегодня
	didTaskToday, err := h.progressRepo.HasUserCompletedTaskToday(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Проверяем засуху (есть ли невыполненные задачи за вчера)
	isDrought, err := h.checkDroughtStatus(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем все задачи пользователя
	tasks, err := h.taskRepo.GetAll(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем все привычки пользователя
	habits, err := h.habitRepo.GetAll(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем все теги пользователя
	tags, err := h.tagRepo.GetByUser(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем грядки (field)
	beds, err := h.bedRepo.GetByUser(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем семена из инвентаря
	inventorySeeds, err := h.userSeedRepo.GetUserSeedsWithDetails(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем товары из магазина
	shopStorage, err := h.goodRepo.GetByUser(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем растения на грядках для seeds (SeedStorage)
	plants, err := h.userPlantRepo.GetWithSeedDetails(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Преобразуем растения в формат IPlant для грядок
	plantMap := make(map[int]IPlant)
	for _, plant := range plants {
		plantMap[plant.BedID] = IPlant{
			ID:           plant.ID,
			Name:         plant.SeedName,
			CurrentStage: plant.CurrentGrowth,
			FinalStage:   plant.TargetGrowth,
			ImgPath:      plant.SeedIcon,
		}
	}

	// Формируем грядки с растениями
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

	// Преобразуем растения в формат SeedStorage для seeds
	seeds := make([]SeedStorage, len(plants))
	for i, plant := range plants {
		seeds[i] = SeedStorage{
			ID:            plant.ID,
			SeedID:        plant.SeedID,
			CurrentGrowth: plant.CurrentGrowth,
			IsWithered:    plant.IsWithered,
			BedID:         plant.BedID,
			CreatedAt:     plant.CreatedAt,
			// Добавляем детали семени
			SeedName:     plant.SeedName,
			TargetGrowth: plant.TargetGrowth,
			Icon:         plant.SeedIcon,
		}
	}

	// Формируем ответ в нужном формате
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

// checkDroughtStatus проверяет статус засухи для пользователя
func (h *UserHandler) checkDroughtStatus(ctx context.Context, userID int64) (bool, error) {
	yesterday := time.Now().AddDate(0, 0, -1)

	// Получаем задачи за вчера
	tasks, err := h.taskRepo.GetByDate(ctx, userID, yesterday)
	if err != nil {
		return false, err
	}

	// Проверяем невыполненные задачи
	for _, task := range tasks {
		if !task.Done {
			return true, nil
		}
	}

	// Проверяем привычки за вчера
	habits, err := h.habitRepo.GetAll(ctx, userID)
	if err != nil {
		return false, err
	}

	for _, habit := range habits {
		if h.shouldHabitBeCompletedYesterday(habit, yesterday) && !habit.Done {
			return true, nil
		}
	}

	return false, nil
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
// @Description Создает нового пользователя или обновляет существующего по MaxID. При создании также создается статистика и начальные грядки.
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
		Gold:                100,
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

	return c.JSON(http.StatusCreated, user)
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
	ID           int    `json:"id"`
	Name         string `json:"name"`
	CurrentStage int    `json:"currentStage"`
	FinalStage   int    `json:"FinalStage"`
	ImgPath      string `json:"imgPath"`
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
	ShopStorage    []model.Good                `json:"shopStorage"`
}

// SeedStorage представляет растение на грядке (для отдельного списка)
type SeedStorage struct {
	ID            int       `json:"id"`
	SeedID        int       `json:"seedId"`
	CurrentGrowth int       `json:"currentGrowth"`
	IsWithered    bool      `json:"isWithered"`
	BedID         int       `json:"bedId"`
	CreatedAt     time.Time `json:"createdAt"`
	// Дополнительные поля из семени
	SeedName     string `json:"seedName,omitempty"`
	TargetGrowth int    `json:"targetGrowth,omitempty"`
	Icon         string `json:"icon,omitempty"`
}
