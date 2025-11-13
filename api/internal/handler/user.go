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
	tagRepo       *repository.TagRepo
	seedRepo      *repository.SeedRepo
	userSeedRepo  *repository.UserSeedRepo
	userPlantRepo *repository.UserPlantRepo
}

func NewUserHandler(
	repo *repository.UserRepo,
	statRepo *repository.UserStatRepo,
	bedRepo *repository.BedRepo,
	taskRepo *repository.TaskRepo,
	tagRepo *repository.TagRepo,
	seedRepo *repository.SeedRepo,
	userSeedRepo *repository.UserSeedRepo,
	userPlantRepo *repository.UserPlantRepo,
) *UserHandler {
	return &UserHandler{
		repo:          repo,
		statRepo:      statRepo,
		bedRepo:       bedRepo,
		taskRepo:      taskRepo,
		tagRepo:       tagRepo,
		seedRepo:      seedRepo,
		userSeedRepo:  userSeedRepo,
		userPlantRepo: userPlantRepo,
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
	hasCompletedToday, err := h.taskRepo.HasCompletedTaskToday(ctx, userID)
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

// StreakInfo представляет информацию о стрике
type StreakInfo struct {
	CurrentStreak int  `json:"current_streak"`
	LongestStreak int  `json:"longest_streak"`
	MissedDay     bool `json:"missed_day"`
}

// DailyStatus представляет статус ежедневной активности
type DailyStatus struct {
	HasCompletedTaskToday bool `json:"has_completed_task_today"`
	MissedDay             bool `json:"missed_day"`
	PlantsWithered        bool `json:"plants_withered"`
	CanRecoverPlants      bool `json:"can_recover_plants"`
}

// SyncUserData godoc
// @Summary Синхронизировать все данные пользователя
// @Description Возвращает все данные пользователя для синхронизации: задачи, теги, статистику, инвентарь, магазин, грядки, растения и информацию о стриках
// @Tags users
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} SyncDataResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/sync [get]
func (h *UserHandler) SyncUserData(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	ctx := context.Background()

	// Получаем базовые данные пользователя
	user, err := h.repo.GetByID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем статистику
	stats, err := h.statRepo.GetByUserID(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем все задачи пользователя
	tasks, err := h.taskRepo.GetAll(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем все теги пользователя
	tags, err := h.tagRepo.GetByUser(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем инвентарь (семена пользователя с деталями)
	inventory, err := h.userSeedRepo.GetUserSeedsWithDetails(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем магазин (все доступные семена)
	userLevel := stats.Level() // Используем вычисляемый уровень
	shop, err := h.userSeedRepo.GetAvailableSeedsForUser(ctx, userID, userLevel)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем грядки
	beds, err := h.bedRepo.GetByUser(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем растения с деталями
	plants, err := h.userPlantRepo.GetWithSeedDetails(ctx, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Проверяем ежедневный статус и обрабатываем логику засыхания растений
	dailyStatus, err := h.processDailyLogic(ctx, userID, user, stats)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Формируем информацию о стрике
	streakInfo := StreakInfo{
		CurrentStreak: stats.CurrentStreak,
		LongestStreak: stats.LongestStreak,
		MissedDay:     dailyStatus.MissedDay,
	}

	// Создаем статистику с вычисляемым уровнем
	statsWithLevel := &UserStatWithLevel{
		UserStat:               stats,
		Level:                  stats.Level(),
		ExperienceForNextLevel: stats.ExperienceForNextLevel(),
	}

	// В response используем statsWithLevel
	response := SyncDataResponse{
		User:        user,
		Stats:       statsWithLevel,
		Tasks:       tasks,
		Tags:        tags,
		Inventory:   inventory,
		Shop:        shop,
		Beds:        beds,
		Plants:      plants,
		StreakInfo:  streakInfo,
		DailyStatus: dailyStatus,
	}

	return c.JSON(http.StatusOK, response)
}

// processDailyLogic обрабатывает ежедневную логику: пропущенные дни и засыхание растений
func (h *UserHandler) processDailyLogic(ctx context.Context, userID int64, user *model.User, stats *model.UserStat) (DailyStatus, error) {
	dailyStatus := DailyStatus{}

	// Получаем последний логин
	lastLogin, err := h.repo.GetLastLogin(ctx, userID)
	if err != nil {
		return dailyStatus, err
	}

	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Проверяем, был ли сегодня выполнен таск
	hasCompletedToday, err := h.taskRepo.HasCompletedTaskToday(ctx, userID)
	if err != nil {
		return dailyStatus, err
	}
	dailyStatus.HasCompletedTaskToday = hasCompletedToday

	// Если пользователь зашел впервые сегодня
	if lastLogin == nil || lastLogin.Before(today) {
		// Обновляем last_login
		user.LastLogin = &now
		if err := h.repo.Update(ctx, user); err != nil {
			return dailyStatus, err
		}

		// Проверяем, пропустил ли пользователь вчера
		if lastLogin != nil {
			yesterday := today.AddDate(0, 0, -1)
			lastLoginDay := time.Date(lastLogin.Year(), lastLogin.Month(), lastLogin.Day(), 0, 0, 0, 0, lastLogin.Location())

			// Если последний логин был позавчера или раньше - пропущен день
			if lastLoginDay.Before(yesterday) {
				dailyStatus.MissedDay = true

				// Помечаем растения как засохшие
				if err := h.userPlantRepo.MarkPlantsAsWithered(ctx, userID); err != nil {
					return dailyStatus, err
				}
				dailyStatus.PlantsWithered = true

				// Сбрасываем стрик
				if err := h.statRepo.ResetStreak(ctx, userID); err != nil {
					return dailyStatus, err
				}
			}
		}
	}

	// Если растения засохли, но пользователь выполнил задание сегодня - можно восстановить
	if dailyStatus.PlantsWithered && dailyStatus.HasCompletedTaskToday {
		dailyStatus.CanRecoverPlants = true
	}

	return dailyStatus, nil
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

// internal/handler/user_stat_handler.go
// GetLevelInfo godoc
// @Summary Получить информацию об уровне
// @Description Возвращает текущий уровень, опыт и прогресс до следующего уровня
// @Tags user-stats
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {object} LevelInfoResponse
// @Router /user-stats/level-info [get]
func (h *UserStatHandler) GetLevelInfo(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	stats, err := h.repo.GetByUserID(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := LevelInfoResponse{
		Level:                  stats.Level(),
		Experience:             stats.Experience,
		ExperienceForNextLevel: stats.ExperienceForNextLevel(),
		ProgressPercent:        float64(stats.Experience) / float64(stats.Level()*100) * 100,
	}

	return c.JSON(http.StatusOK, response)
}

// LevelInfoResponse представляет информацию об уровне
type LevelInfoResponse struct {
	Level                  int     `json:"level"`
	Experience             int64   `json:"experience"`
	ExperienceForNextLevel int64   `json:"experience_for_next_level"`
	ProgressPercent        float64 `json:"progress_percent"`
}

// DTO для запросов

// UserCreateRequest представляет запрос на создание/обновление пользователя
type UserCreateRequest struct {
	MaxID    int64  `json:"max_id" example:"123456789"`
	Username string `json:"username" example:"john_doe"`
}

// UserUpdateRequest представляет запрос на обновление пользователя
type UserUpdateRequest struct {
	Username string `json:"username" example:"john_doe_updated"`
}

// UserStatWithLevel представляет статистику с вычисляемым уровнем
type UserStatWithLevel struct {
	*model.UserStat
	Level                  int   `json:"level"`
	ExperienceForNextLevel int64 `json:"experience_for_next_level"`
}

// SyncDataResponse обновите:
type SyncDataResponse struct {
	User        *model.User                 `json:"user"`
	Stats       *UserStatWithLevel          `json:"stats"`
	Tasks       []model.Task                `json:"tasks"`
	Tags        []model.Tag                 `json:"tags"`
	Inventory   []model.UserSeedWithDetails `json:"inventory"`
	Shop        []model.SeedWithUserData    `json:"shop"`
	Beds        []model.Bed                 `json:"beds"`
	Plants      []model.UserPlantWithSeed   `json:"plants"`
	StreakInfo  StreakInfo                  `json:"streak_info"`
	DailyStatus DailyStatus                 `json:"daily_status"`
}
