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

type CategoryHandler struct {
	repo *repository.CategoryRepo
}

func NewCategoryHandler(repo *repository.CategoryRepo) *CategoryHandler {
	return &CategoryHandler{repo: repo}
}

func (h *CategoryHandler) Create(c echo.Context) error {
	var cat model.Category
	if err := c.Bind(&cat); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}
	cat.UserID = 1    // TODO
	cat.CreatedBy = 1 // TODO
	cat.UpdatedBy = 1 // TODO
	if err := h.repo.Create(context.Background(), &cat); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, cat)
}

func (h *CategoryHandler) GetAll(c echo.Context) error {
	userID := 1 // TODO
	cats, err := h.repo.GetAll(context.Background(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, cats)
}

func (h *CategoryHandler) Update(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	var cat model.Category
	if err := c.Bind(&cat); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	cat.ID = id
	cat.UserID = 1    // TODO
	cat.UpdatedBy = 1 // TODO

	if err := h.repo.Update(context.Background(), &cat); err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, cat)
}

func (h *CategoryHandler) Delete(c echo.Context) error {
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
