package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
	"github.com/RinatHar/FarmFocus/api/internal/utils"
	"github.com/labstack/echo/v4"
)

type HabitHandler struct {
	BaseHandler
	repo          *repository.HabitRepo
	progressRepo  *repository.ProgressLogRepo
	userStatRepo  *repository.UserStatRepo
	userPlantRepo *repository.UserPlantRepo
}

func NewHabitHandler(
	repo *repository.HabitRepo,
	progressRepo *repository.ProgressLogRepo,
	userStatRepo *repository.UserStatRepo,
	userPlantRepo *repository.UserPlantRepo) *HabitHandler {
	return &HabitHandler{
		repo:          repo,
		progressRepo:  progressRepo,
		userStatRepo:  userStatRepo,
		userPlantRepo: userPlantRepo,
	}
}

// Create godoc
// @Summary Создать новую привычку
// @Description Создает новую привычку для текущего пользователя с базовым количеством опыта с учетом count
// @Tags habits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param request body HabitCreateRequest true "Данные для создания привычки"
// @Success 200 {object} model.Habit
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /habits [post]
func (h *HabitHandler) Create(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req HabitCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Устанавливаем базовое количество опыта с учетом сложности и count
	baseXP := utils.GetBaseXPForHabit(req.Difficulty, req.Count)

	habit := model.Habit{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		Done:        false,
		Count:       req.Count,
		Period:      req.Period,
		Every:       req.Every,
		StartDate:   req.StartDate,
		XPReward:    baseXP,
		TagID:       req.TagID,
		CreatedAt:   time.Now(),
	}

	if err := h.repo.Create(context.Background(), &habit); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, habit)
}

// GetAll godoc
// @Summary Получить все привычки пользователя
// @Description Возвращает список всех привычек текущего пользователя
// @Tags habits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param done query boolean false "Фильтр по статусу выполнения"
// @Param period query string false "Фильтр по периоду (day/week/month)"
// @Param tag_id query integer false "Фильтр по ID тега"
// @Success 200 {array} model.Habit
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /habits [get]
func (h *HabitHandler) GetAll(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	// Обработка query параметров
	doneStr := c.QueryParam("done")
	period := c.QueryParam("period")
	tagIDStr := c.QueryParam("tag_id")

	var habits []model.Habit

	if tagIDStr != "" {
		tagID, err := strconv.Atoi(tagIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tag ID"})
		}
		habits, err = h.repo.GetByTag(context.Background(), userID, tagID)
	} else if period != "" {
		habits, err = h.repo.GetByPeriod(context.Background(), userID, period)
	} else if doneStr != "" {
		done, err := strconv.ParseBool(doneStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid done value"})
		}
		habits, err = h.repo.GetByStatus(context.Background(), userID, done)
	} else {
		habits, err = h.repo.GetAll(context.Background(), userID)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, habits)
}

// GetByID godoc
// @Summary Получить привычку по ID
// @Description Возвращает привычку по указанному ID
// @Tags habits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Habit ID"
// @Success 200 {object} model.Habit
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /habits/{id} [get]
func (h *HabitHandler) GetByID(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid habit ID"})
	}

	habit, err := h.repo.GetByID(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, habit)
}

// Update godoc
// @Summary Обновить привычку
// @Description Обновляет информацию о привычке
// @Tags habits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Habit ID"
// @Param request body HabitUpdateRequest true "Обновленные данные привычки"
// @Success 200 {object} model.Habit
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /habits/{id} [put]
func (h *HabitHandler) Update(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid habit ID"})
	}

	var req HabitUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	habit := model.Habit{
		ID:          id,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		Done:        req.Done,
		Count:       req.Count,
		Period:      req.Period,
		Every:       req.Every,
		StartDate:   req.StartDate,
		XPReward:    req.XPReward,
		TagID:       req.TagID,
	}

	if err := h.repo.Update(context.Background(), &habit); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, habit)
}

// Delete godoc
// @Summary Удалить привычку
// @Description Удаляет привычку по ID, предварительно удаляя связанные записи в логе прогресса
// @Tags habits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Habit ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /habits/{id} [delete]
func (h *HabitHandler) Delete(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid habit ID"})
	}

	// Проверяем существование привычки и принадлежность пользователю
	if _, err := h.repo.GetByID(context.Background(), id, userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Удаляем связанные записи в progress_log
	if err := h.progressRepo.DeleteByHabitID(context.Background(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Удаляем привычку
	err = h.repo.Delete(context.Background(), id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// MarkAsDone godoc
// @Summary Пометить привычку как выполненную
// @Description Помечает привычку как выполненную, увеличивает счетчик, добавляет опыт пользователю по формуле, создает запись в логе прогресса, увеличивает рост всех активных растений пользователя на 1 и увеличивает стрик если это первое выполнение сегодня
// @Tags habits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Habit ID"
// @Success 200 {object} HabitCompletionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /habits/{id}/done [patch]
func (h *HabitHandler) MarkAsDone(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid habit ID"})
	}

	// Получаем привычку
	habit, err := h.repo.GetByID(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Проверяем, не выполнена ли уже привычка
	if habit.Done {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Habit is already marked as done",
		})
	}

	// Получаем статистику пользователя для расчета уровня
	stats, err := h.userStatRepo.GetByUserID(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Рассчитываем опыт по формуле: базовый опыт * (коэф уровня + коэф сложности)
	userLevel := stats.Level()
	calculatedXP := utils.CalculateTaskXP(habit.XPReward, userLevel, habit.Difficulty)

	// Проверяем, выполнены ли сегодня задачи или привычки
	hasCompletedTaskToday, err := h.progressRepo.HasUserCompletedTaskToday(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	hasCompletedHabitToday, err := h.progressRepo.HasUserCompletedHabitToday(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Если сегодня еще не выполнено ни задач, ни привычек - увеличиваем стрик
	shouldIncrementStreak := !hasCompletedTaskToday && !hasCompletedHabitToday

	// Если это первое выполнение сегодня и пользователь в засухе - сбрасываем засуху
	if shouldIncrementStreak && stats.IsDrought {
		h.userStatRepo.ResetDrought(context.Background(), userID)
	}

	// Помечаем привычку как выполненную и увеличиваем счетчик
	err = h.repo.MarkAsDoneAndIncrementCount(context.Background(), id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Увеличиваем стрик если это первое выполнение сегодня
	if shouldIncrementStreak {
		if err := h.userStatRepo.IncrementStreak(context.Background(), userID); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	// Создаем запись в progress_log
	progressLog := &model.ProgressLog{
		UserID:     userID,
		HabitID:    &id,
		XPEarned:   calculatedXP,
		GoldEarned: 0,
		CreatedAt:  time.Now(),
	}

	if err := h.progressRepo.Create(context.Background(), progressLog); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Добавляем опыт пользователю
	if calculatedXP > 0 {
		if err := h.userStatRepo.AddExperience(context.Background(), userID, int64(calculatedXP)); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	// Увеличиваем рост всех активных растений пользователя на 1
	plantsGrown, err := h.growUserPlants(context.Background(), userID, 1)
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение основной операции
	}

	return c.JSON(http.StatusOK, HabitCompletionResponse{
		XPEarned:    calculatedXP,
		PlantsGrown: plantsGrown,
	})
}

// MarkAsUndone godoc
// @Summary Пометить привычку как невыполненную
// @Description Помечает привычку как невыполненную, уменьшает счетчик и возвращает опыт обратно на основе последней записи в логе
// @Tags habits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Habit ID"
// @Success 200 {object} HabitUndoResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /habits/{id}/undone [patch]
func (h *HabitHandler) MarkAsUndone(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid habit ID"})
	}

	// Получаем привычку для проверки существования и текущего состояния
	habit, err := h.repo.GetByID(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Проверяем, выполнена ли привычка
	if !habit.Done {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Habit is not marked as done",
		})
	}

	// Проверяем, что счетчик больше 0
	if habit.Count <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Habit count is already 0",
		})
	}

	// Ищем последнюю запись в progress_log для этой привычки
	lastProgress, err := h.progressRepo.GetLastByHabitID(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "No completion record found for this habit",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Проверяем, что запись содержит положительный опыт (привычка была выполнена)
	if lastProgress.XPEarned <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Habit was not completed or already undone",
		})
	}

	xpToReturn := lastProgress.XPEarned

	// Помечаем привычку как невыполненную и уменьшаем счетчик
	err = h.repo.MarkAsUndoneAndDecrementCount(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Возвращаем опыт (вычитаем)
	if xpToReturn > 0 {
		if err := h.userStatRepo.RemoveExperience(context.Background(), userID, int64(xpToReturn)); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	// Создаем запись в progress_log о возврате опыта
	progressLog := &model.ProgressLog{
		UserID:     userID,
		HabitID:    &id,
		XPEarned:   -xpToReturn, // Отрицательное значение для возврата
		GoldEarned: 0,
		CreatedAt:  time.Now(),
	}

	if err := h.progressRepo.Create(context.Background(), progressLog); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, HabitUndoResponse{
		XPEarned: -xpToReturn,
	})
}

// IncrementCount godoc
// @Summary Увеличить счетчик привычки
// @Description Увеличивает счетчик привычки на 1
// @Tags habits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Habit ID"
// @Success 200 {object} HabitIncrementResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /habits/{id}/increment [patch]
func (h *HabitHandler) IncrementCount(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid habit ID"})
	}

	// Увеличиваем счетчик
	err = h.repo.IncrementCount(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Получаем обновленную привычку для ответа
	habit, err := h.repo.GetByID(context.Background(), id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, HabitIncrementResponse{
		Count: habit.Count,
	})
}

// ResetCount godoc
// @Summary Сбросить счетчик привычки
// @Description Сбрасывает счетчик привычки до 0
// @Tags habits
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Habit ID"
// @Success 200 {object} HabitResetResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /habits/{id}/reset [patch]
func (h *HabitHandler) ResetCount(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid habit ID"})
	}

	// Сбрасываем счетчик
	err = h.repo.ResetCount(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, HabitResetResponse{
		Message: "Habit count reset successfully",
	})
}

// growUserPlants увеличивает рост всех активных растений пользователя
func (h *HabitHandler) growUserPlants(ctx context.Context, userID int64, growthAmount int) (int, error) {
	userPlants, err := h.userPlantRepo.GetByUser(ctx, userID)
	if err != nil {
		return 0, err
	}

	grownCount := 0
	for _, plant := range userPlants {
		_, err := h.userPlantRepo.AddGrowth(ctx, plant.ID, growthAmount)
		if err != nil {
			continue
		}
		grownCount++
	}

	return grownCount, nil
}

// DTO для запросов

// HabitCompletionResponse представляет ответ при завершении привычки
type HabitCompletionResponse struct {
	XPEarned    int `json:"xpEarned" example:"150"`
	PlantsGrown int `json:"plantsGrown" example:"3"`
}

// HabitUndoResponse представляет ответ при отмене выполнения привычки
type HabitUndoResponse struct {
	XPEarned int `json:"xpEarned" example:"150"`
}

// HabitIncrementResponse представляет ответ при увеличении счетчика
type HabitIncrementResponse struct {
	Count int `json:"count" example:"5"`
}

// HabitResetResponse представляет ответ при сбросе счетчика
type HabitResetResponse struct {
	Message string `json:"message" example:"Habit count reset successfully"`
}

// HabitCreateRequest представляет запрос на создание привычки
type HabitCreateRequest struct {
	Title       string    `json:"title" example:"Утренняя зарядка"`
	Description *string   `json:"description,omitempty" example:"Ежедневная утренняя зарядка 15 минут"`
	Difficulty  string    `json:"difficulty" example:"medium"`
	Count       int       `json:"count" example:"0"`
	Period      string    `json:"period" example:"day"`
	Every       int       `json:"every" example:"1"`
	StartDate   time.Time `json:"startDate" example:"2024-01-15T00:00:00Z"`
	XPReward    int       `json:"xpReward" example:"50"`
	TagID       *int      `json:"tagId,omitempty" example:"1"`
}

// HabitUpdateRequest представляет запрос на обновление привычки
type HabitUpdateRequest struct {
	Title       string    `json:"title" example:"Утренняя зарядка - обновлено"`
	Description *string   `json:"description,omitempty" example:"Обновленное описание привычки"`
	Difficulty  string    `json:"difficulty" example:"hard"`
	Done        bool      `json:"done" example:"false"`
	Count       int       `json:"count" example:"5"`
	Period      string    `json:"period" example:"week"`
	Every       int       `json:"every" example:"3"`
	StartDate   time.Time `json:"startDate" example:"2024-01-20T00:00:00Z"`
	XPReward    int       `json:"xpReward" example:"75"`
	TagID       *int      `json:"tagId,omitempty" example:"2"`
}
