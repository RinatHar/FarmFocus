package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/RinatHar/FarmFocus/api/internal/model"
	"github.com/RinatHar/FarmFocus/api/internal/repository"
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
// @Description Создает новую привычку для текущего пользователя
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
		XPReward:    req.XPReward,
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
// @Description Удаляет привычку по ID
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

	err = h.repo.Delete(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// MarkAsDone godoc
// @Summary Пометить привычку как выполненную
// @Description Помечает привычку как выполненную, добавляет опыт пользователю, создает запись в логе прогресса и увеличивает рост всех активных растений пользователя на 1
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

	// Получаем привычку чтобы узнать награду
	habit, err := h.repo.GetByID(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Помечаем привычку как выполненную
	err = h.repo.MarkAsDone(context.Background(), id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Создаем запись в progress_log
	progressLog := &model.ProgressLog{
		UserID:     userID,
		HabitID:    &id,
		XPEarned:   habit.XPReward,
		GoldEarned: 0,
		CreatedAt:  time.Now(),
	}

	if err := h.progressRepo.Create(context.Background(), progressLog); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Добавляем опыт пользователю
	if habit.XPReward > 0 {
		if err := h.userStatRepo.AddExperience(context.Background(), userID, int64(habit.XPReward)); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	// Увеличиваем рост всех активных растений пользователя на 1
	plantsGrown, err := h.growUserPlants(context.Background(), userID, 1)
	if err != nil {
		// Логируем ошибку, но не прерываем выполнение основной операции
	}

	return c.JSON(http.StatusOK, HabitCompletionResponse{
		Message:     "Habit completed successfully",
		XPEarned:    habit.XPReward,
		HabitID:     id,
		PlantsGrown: plantsGrown,
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
// @Success 200 {object} map[string]string
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

	err = h.repo.IncrementCount(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Habit count incremented"})
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
// @Success 200 {object} map[string]string
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

	err = h.repo.ResetCount(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "Habit count reset"})
}

// MarkAsUndone godoc
// @Summary Пометить привычку как невыполненную
// @Description Помечает привычку как невыполненную (сбрасывает статус)
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

	err = h.repo.MarkAsUndone(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
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

// HabitCompletionResponse представляет ответ при завершении привычки
type HabitCompletionResponse struct {
	Message     string `json:"message" example:"Habit completed successfully"`
	XPEarned    int    `json:"xpEarned" example:"50"`
	HabitID     int    `json:"habitId" example:"123"`
	PlantsGrown int    `json:"plantsGrown" example:"3"`
}
