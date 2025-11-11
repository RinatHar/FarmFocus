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

type TaskHandler struct {
	BaseHandler
	repo         *repository.TaskRepo
	progressRepo *repository.ProgressLogRepo
	userStatRepo *repository.UserStatRepo
}

func NewTaskHandler(repo *repository.TaskRepo, progressRepo *repository.ProgressLogRepo, userStatRepo *repository.UserStatRepo) *TaskHandler {
	return &TaskHandler{
		repo:         repo,
		progressRepo: progressRepo,
		userStatRepo: userStatRepo,
	}
}

// Create godoc
// @Summary Создать новую задачу
// @Description Создает новую задачу для текущего пользователя
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

	task := model.Task{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		DueDate:     req.DueDate,
		IsDone:      false,
		XPReward:    req.XPReward,
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
// @Success 200 {array} model.Task
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks [get]
func (h *TaskHandler) GetAll(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	tasks, err := h.repo.GetAll(context.Background(), userID)
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
		DueDate:     req.DueDate,
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
// @Description Удаляет задачу по ID
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

	err = h.repo.Delete(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetByStatus godoc
// @Summary Получить задачи по статусу
// @Description Возвращает задачи по статусу выполнения (pending или completed)
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param status query string true "Статус задач (pending/completed)"
// @Success 200 {array} model.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/status [get]
func (h *TaskHandler) GetByStatus(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	status := c.QueryParam("status")
	if status == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "status parameter is required"})
	}

	isDone := false
	if status == "completed" {
		isDone = true
	} else if status != "pending" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "status must be 'pending' or 'completed'"})
	}

	tasks, err := h.repo.GetByStatus(context.Background(), userID, isDone)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tasks)
}

// MarkAsDone godoc
// @Summary Пометить задачу как выполненную
// @Description Помечает задачу как выполненную, добавляет опыт пользователю и создает запись в логе прогресса
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

	// Получаем задачу чтобы узнать награду
	task, err := h.repo.GetByID(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Помечаем задачу как выполненную
	err = h.repo.MarkAsDone(context.Background(), id, userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Создаем запись в progress_log
	progressLog := &model.ProgressLog{
		UserID:     userID,
		TaskID:     id,
		XPEarned:   task.XPReward,
		GoldEarned: 0, // можно добавить награду золотом если нужно
		CreatedAt:  time.Now(),
	}

	if err := h.progressRepo.Create(context.Background(), progressLog); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Добавляем опыт пользователю
	if task.XPReward > 0 {
		if err := h.userStatRepo.AddExperience(context.Background(), userID, int64(task.XPReward)); err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
	}

	return c.JSON(http.StatusOK, TaskCompletionResponse{
		Message:  "Task completed successfully",
		XPEarned: task.XPReward,
		TaskID:   id,
	})
}

// MarkAsUndone godoc
// @Summary Пометить задачу как невыполненную
// @Description Помечает задачу как невыполненную (сбрасывает статус)
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

	err = h.repo.MarkAsUndone(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetOverdue godoc
// @Summary Получить просроченные задачи
// @Description Возвращает список просроченных задач (с истекшим сроком выполнения)
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.Task
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/overdue [get]
func (h *TaskHandler) GetOverdue(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	tasks, err := h.repo.GetOverdue(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tasks)
}

// GetByTag godoc
// @Summary Получить задачи по тегу
// @Description Возвращает задачи, привязанные к определенному тегу
// @Tags tasks
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param tagId path int true "Tag ID"
// @Success 200 {array} model.Task
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tasks/tag/{tagId} [get]
func (h *TaskHandler) GetByTag(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	tagID, err := strconv.Atoi(c.Param("tagId"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tag ID"})
	}

	tasks, err := h.repo.GetByTag(context.Background(), userID, tagID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tasks)
}

// DTO для запросов

// TaskCreateRequest представляет запрос на создание задачи
type TaskCreateRequest struct {
	Title       string     `json:"title" example:"Завершить проект"`
	Description string     `json:"description" example:"Нужно завершить все задачи по проекту до конца недели"`
	DueDate     *time.Time `json:"due_date,omitempty" example:"2024-01-15T00:00:00Z"`
	XPReward    int        `json:"xp_reward" example:"100"`
	TagID       *int       `json:"tag_id,omitempty" example:"1"`
}

// TaskUpdateRequest представляет запрос на обновление задачи
type TaskUpdateRequest struct {
	Title       string     `json:"title" example:"Завершить проект - обновлено"`
	Description string     `json:"description" example:"Обновленное описание задачи"`
	DueDate     *time.Time `json:"due_date,omitempty" example:"2024-01-20T00:00:00Z"`
	XPReward    int        `json:"xp_reward" example:"150"`
	TagID       *int       `json:"tag_id,omitempty" example:"2"`
}

// TaskCompletionResponse представляет ответ при завершении задачи
type TaskCompletionResponse struct {
	Message  string `json:"message" example:"Task completed successfully"`
	XPEarned int    `json:"xp_earned" example:"100"`
	TaskID   int    `json:"task_id" example:"123"`
}
