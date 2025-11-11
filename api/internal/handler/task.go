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

// Create создает новую задачу
func (h *TaskHandler) Create(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var task model.Task
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	task.UserID = userID
	task.CreatedAt = time.Now()

	if err := h.repo.Create(context.Background(), &task); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, task)
}

// GetAll возвращает все задачи пользователя
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

// GetByID возвращает задачу по ID
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

// Update обновляет задачу
func (h *TaskHandler) Update(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid task ID"})
	}

	var task model.Task
	if err := c.Bind(&task); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	task.ID = id
	task.UserID = userID

	if err := h.repo.Update(context.Background(), &task); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, task)
}

// Delete удаляет задачу
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

// GetByStatus возвращает задачи по статусу выполнения
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

// MarkAsDone помечает задачу как выполненную и создает запись в progress_log
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

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message":   "Task completed successfully",
		"xp_earned": task.XPReward,
		"task_id":   id,
	})
}

// MarkAsUndone помечает задачу как невыполненную
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

// GetOverdue возвращает просроченные задачи
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

// GetByTag возвращает задачи по тегу
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
