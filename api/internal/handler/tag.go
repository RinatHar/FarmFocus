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
	repo      *repository.TagRepo
	taskRepo  *repository.TaskRepo
	habitRepo *repository.HabitRepo
}

func NewTagHandler(
	repo *repository.TagRepo,
	taskRepo *repository.TaskRepo,
	habitRepo *repository.HabitRepo,
) *TagHandler {
	return &TagHandler{
		repo:      repo,
		taskRepo:  taskRepo,
		habitRepo: habitRepo,
	}
}

// GetAll godoc
// @Summary Получить все теги пользователя
// @Description Возвращает список всех тегов текущего пользователя
// @Tags tags
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.Tag
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tags [get]
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

// GetByID godoc
// @Summary Получить тег по ID
// @Description Возвращает тег по указанному ID
// @Tags tags
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Tag ID"
// @Success 200 {object} model.Tag
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tags/{id} [get]
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

// Create godoc
// @Summary Создать новый тег
// @Description Создает новый тег для текущего пользователя
// @Tags tags
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param request body TagCreateRequest true "Данные для создания тега"
// @Success 200 {object} model.Tag
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tags [post]
func (h *TagHandler) Create(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	var req TagCreateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tag := model.Tag{
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	}

	if err := h.repo.Create(context.Background(), &tag); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, tag)
}

// Update godoc
// @Summary Обновить тег
// @Description Обновляет информацию о теге
// @Tags tags
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Tag ID"
// @Param request body TagUpdateRequest true "Обновленные данные тега"
// @Success 200 {object} model.Tag
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tags/{id} [put]
func (h *TagHandler) Update(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tag ID"})
	}

	var req TagUpdateRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	tag := model.Tag{
		ID:     id,
		UserID: userID,
		Name:   req.Name,
		Color:  req.Color,
	}

	if err := h.repo.Update(context.Background(), &tag); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, tag)
}

// Delete godoc
// @Summary Удалить тег
// @Description Удаляет тег и сбрасывает его у всех связанных задач и привычек
// @Tags tags
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Param id path int true "Tag ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 403 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tags/{id} [delete]
func (h *TagHandler) Delete(c echo.Context) error {
	userID, err := h.GetUserIDFromContext(c)
	if err != nil {
		return err
	}

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid tag ID"})
	}

	ctx := context.Background()

	// Сначала сбрасываем тег у всех задач пользователя
	if err := h.taskRepo.ResetTag(ctx, userID, id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to reset tag in tasks: " + err.Error()})
	}

	// Сбрасываем тег у всех привычек пользователя
	if err := h.habitRepo.ResetTag(ctx, userID, id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to reset tag in habits: " + err.Error()})
	}

	// Теперь удаляем сам тег
	err = h.repo.Delete(ctx, id, userID)
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

// GetWithTaskCount godoc
// @Summary Получить теги с количеством задач
// @Description Возвращает теги с информацией о количестве привязанных задач
// @Tags tags
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param X-User-ID header string true "User ID"
// @Success 200 {array} model.TagWithCount
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /tags/with-count [get]
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

// DTO для запросов

// TagCreateRequest представляет запрос на создание тега
type TagCreateRequest struct {
	Name  string `json:"name" example:"Работа"`
	Color string `json:"color" example:"#FF5733"`
}

// TagUpdateRequest представляет запрос на обновление тега
type TagUpdateRequest struct {
	Name  string `json:"name" example:"Работа - обновлено"`
	Color string `json:"color" example:"#33FF57"`
}
