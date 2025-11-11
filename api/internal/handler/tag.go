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

type TagHandler struct {
	BaseHandler
	repo *repository.TagRepo
}

func NewTagHandler(repo *repository.TagRepo) *TagHandler {
	return &TagHandler{repo: repo}
}

// GetAll возвращает все теги пользователя
func (h *TagHandler) GetAll(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	tags, err := h.repo.GetByUser(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tags)
}

// GetByID возвращает тег по ID
func (h *TagHandler) GetByID(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tag ID"})
	}

	tag, err := h.repo.GetByID(context.Background(), id)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Проверяем, что тег принадлежит пользователю
	if tag.UserID != userID {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "access denied"})
	}

	return c.JSON(http.StatusOK, tag)
}

// Create создает новый тег
func (h *TagHandler) Create(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var tag model.Tag
	if err := c.Bind(&tag); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tag.UserID = userID

	if err := h.repo.Create(context.Background(), &tag); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tag)
}

// Update обновляет тег
func (h *TagHandler) Update(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tag ID"})
	}

	var tag model.Tag
	if err := c.Bind(&tag); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tag.ID = id
	tag.UserID = userID

	if err := h.repo.Update(context.Background(), &tag); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tag)
}

// Delete удаляет тег
func (h *TagHandler) Delete(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tag ID"})
	}

	err = h.repo.Delete(context.Background(), id, userID)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		if strings.Contains(err.Error(), "cannot delete tag") {
			return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

// GetWithTaskCount возвращает теги с количеством задач
func (h *TagHandler) GetWithTaskCount(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	tags, err := h.repo.GetWithTaskCount(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tags)
}
