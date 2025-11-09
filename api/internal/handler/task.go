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

type TaskHandler struct {
	repo *repository.TaskRepo
}

func NewTaskHandler(repo *repository.TaskRepo) *TaskHandler {
	return &TaskHandler{repo: repo}
}

func (h *TaskHandler) Create(c echo.Context) error {
	var t model.Task
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	t.UserID = 1    // TODO
	t.CreatedBy = 1 // TODO
	t.UpdatedBy = 1 // TODO
	if err := h.repo.Create(context.Background(), &t); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, t)
}

func (h *TaskHandler) GetAll(c echo.Context) error {
	userID := 1 // TODO
	tasks, err := h.repo.GetAll(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tasks)
}

func (h *TaskHandler) GetByID(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := 1 // TODO
	task, err := h.repo.GetByID(context.Background(), id, userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "task not found"})
	}
	return c.JSON(http.StatusOK, task)
}

func (h *TaskHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var t model.Task
	if err := c.Bind(&t); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	t.ID = id
	t.UserID = 1    // TODO
	t.UpdatedBy = 1 // TODO

	if err := h.repo.Update(context.Background(), &t); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, t)
}

func (h *TaskHandler) Delete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	userID := 1    // TODO
	deletedBy := 1 // TODO

	err := h.repo.Delete(context.Background(), id, userID, deletedBy)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
