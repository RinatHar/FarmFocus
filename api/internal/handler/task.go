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

type TaskHandler struct {
	BaseHandler
	repo          *repository.TaskRepo
	progressRepo  *repository.ProgressLogRepo
	userStatRepo  *repository.UserStatRepo
	userPlantRepo *repository.UserPlantRepo
}

func NewTaskHandler(
	repo *repository.TaskRepo,
	progressRepo *repository.ProgressLogRepo,
	userStatRepo *repository.UserStatRepo,
	userPlantRepo *repository.UserPlantRepo) *TaskHandler {
	return &TaskHandler{
		repo:          repo,
		progressRepo:  progressRepo,
		userStatRepo:  userStatRepo,
		userPlantRepo: userPlantRepo,
	}
}

// Create godoc
// @Summary Создать новую задачу
// @Description Создает новую задачу для текущего пользователя с базовым количеством опыта
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param request body TaskCreateRequest true "Данные для создания задачи"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [post]
func (h *TaskHandler) Create(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req TaskCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	// Устанавливаем базовое количество опыта в зависимости от сложности
	calcXP := utils.GetBaseXPForDifficulty(req.Difficulty)

	task := model.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		Date:        req.Date,
		Done:        false,
		XPReward:    calcXP,
		TagID:       req.TagID,
		CreatedAt:   time.Now(),
	}

	if err := h.repo.Create(context.Background(), &task); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, task)
}

// GetAll godoc
// @Summary Получить все задачи пользователя
// @Description Возвращает список всех задач текущего пользователя
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param done query boolean false "Фильтр по статусу выполнения"
// @Param date query string false "Фильтр по дате (формат: 2006-01-02)"
// @Param tag_id query integer false "Фильтр по ID тега"
// @Success 200 {array} model.Task
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
func (h *TaskHandler) GetAll(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	// Обработка query параметров
	doneStr := c.QueryParam("done")
	dateStr := c.QueryParam("date")
	tagIDStr := c.QueryParam("tag_id")

	var tasks []model.Task

	if tagIDStr != "" {
		tagID, err := strconv.Atoi(tagIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tag ID"})
		}
		tasks, err = h.repo.GetByTag(context.Background(), userID, tagID)
	} else if dateStr != "" {
		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid date format, use YYYY-MM-DD"})
		}
		tasks, err = h.repo.GetByDate(context.Background(), userID, date)
	} else if doneStr != "" {
		done, err := strconv.ParseBool(doneStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid done value"})
		}
		tasks, err = h.repo.GetByStatus(context.Background(), userID, done)
	} else {
		tasks, err = h.repo.GetAll(context.Background(), userID)
	}

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tasks)
}

// GetByID godoc
// @Summary Получить задачу по ID
// @Description Возвращает задачу по указанному ID
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Task ID"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [get]
func (h *TaskHandler) GetByID(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task ID"})
	}

	task, err := h.repo.GetByID(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, task)
}

// Update godoc
// @Summary Обновить задачу
// @Description Обновляет информацию о задаче
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Task ID"
// @Param request body TaskUpdateRequest true "Обновленные данные задачи"
// @Success 200 {object} model.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [put]
func (h *TaskHandler) Update(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task ID"})
	}

	var req TaskUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	task := model.Task{
		ID:          id,
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		Date:        req.Date,
		Done:        req.Done,
		XPReward:    req.XPReward,
		TagID:       req.TagID,
	}

	if err := h.repo.Update(context.Background(), &task); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, task)
}

// Delete godoc
// @Summary Удалить задачу
// @Description Удаляет задачу по ID, предварительно удаляя связанные записи в логе прогресса
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Task ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id} [delete]
func (h *TaskHandler) Delete(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task ID"})
	}

	// Проверяем существование задачи и принадлежность пользователю
	if _, err := h.repo.GetByID(context.Background(), id, userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Удаляем связанные записи в progress_log
	if err := h.progressRepo.DeleteByTaskID(context.Background(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Удаляем задачу
	err = h.repo.Delete(context.Background(), id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// growUserPlants увеличивает рост всех активных растений пользователя
func (h *TaskHandler) growUserPlants(ctx context.Context, userID int64, growthAmount int) (int, error) {
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

// MarkAsDone godoc
// @Summary Пометить задачу как выполненную
// @Description Помечает задачу как выполненную, добавляет опыт пользователю по формуле, создает запись в логе прогресса, увеличивает рост всех активных растений пользователя на 1 и увеличивает стрик если это первое выполнение сегодня
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Task ID"
// @Success 200 {object} TaskCompletionResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id}/done [patch]
func (h *TaskHandler) MarkAsDone(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task ID"})
	}

	// Получаем задачу
	task, err := h.repo.GetByID(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Проверяем, не выполнена ли уже задача
	if task.Done {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Task is already marked as done",
		})
	}

	// Получаем статистику пользователя для расчета уровня
	stats, err := h.userStatRepo.GetByUserID(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Рассчитываем опыт по формуле: базовый опыт * (коэф уровня + коэф сложности)
	userLevel := stats.Level()
	calculatedXP := utils.CalculateTaskXP(task.XPReward, userLevel, task.Difficulty)

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

	// Помечаем задачу как выполненную
	err = h.repo.MarkAsDone(context.Background(), id, userID)
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
		TaskID:     &id,
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

	return c.JSON(http.StatusOK, TaskCompletionResponse{
		XPEarned:    calculatedXP,
		PlantsGrown: plantsGrown,
	})
}

// MarkAsUndone godoc
// @Summary Пометить задачу как невыполненную
// @Description Помечает задачу как невыполненную и возвращает опыт обратно на основе последней записи в логе
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Task ID"
// @Success 200 {object} TaskUndoResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/{id}/undone [patch]
func (h *TaskHandler) MarkAsUndone(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task ID"})
	}

	// Получаем задачу для проверки существования

	if _, err := h.repo.GetByID(context.Background(), id, userID); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Ищем последнюю запись в progress_log для этой задачи
	lastProgress, err := h.progressRepo.GetLastByTaskID(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "No completion record found for this task",
			})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Проверяем, что запись содержит положительный опыт (задача была выполнена)
	if lastProgress.XPEarned <= 0 {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Task was not completed or already undone",
		})
	}

	xpToReturn := lastProgress.XPEarned

	// Помечаем задачу как невыполненную
	err = h.repo.MarkAsUndone(context.Background(), id, userID)
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
		TaskID:     &id,
		XPEarned:   -xpToReturn, // Отрицательное значение для возврата
		GoldEarned: 0,
		CreatedAt:  time.Now(),
	}

	if err := h.progressRepo.Create(context.Background(), progressLog); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, TaskUndoResponse{
		XPEarned: -xpToReturn,
	})
}

// DTO для запросов

// TaskCompletionResponse представляет ответ при завершении задачи
type TaskCompletionResponse struct {
	XPEarned    int `json:"xpEarned" example:"150"`
	PlantsGrown int `json:"plantsGrown" example:"3"`
}

// TaskUndoResponse представляет ответ при отмене выполнения задачи
type TaskUndoResponse struct {
	XPEarned int `json:"xpEarned" example:"150"`
}

// TaskCreateRequest представляет запрос на создание задачи
type TaskCreateRequest struct {
	Title       string     `json:"title" example:"Завершить проект"`
	Description *string    `json:"description,omitempty" example:"Нужно завершить все задачи по проекту до конца недели"`
	Difficulty  string     `json:"difficulty" example:"medium"`
	Date        *time.Time `json:"date,omitempty" example:"2024-01-15T00:00:00Z"`
	TagID       *int       `json:"tagId,omitempty" example:"1"`
}

// TaskUpdateRequest представляет запрос на обновление задачи
type TaskUpdateRequest struct {
	Title       string     `json:"title" example:"Завершить проект - обновлено"`
	Description *string    `json:"description,omitempty" example:"Обновленное описание задачи"`
	Difficulty  string     `json:"difficulty" example:"hard"`
	Date        *time.Time `json:"date,omitempty" example:"2024-01-20T00:00:00Z"`
	Done        bool       `json:"done" example:"false"`
	XPReward    int        `json:"xpReward" example:"150"`
	TagID       *int       `json:"tagId,omitempty" example:"2"`
}
