package handler

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/burbble/marketplace/internal/service"
)

type CategoryHandler struct {
	svc service.CategoryService
}

func NewCategoryHandler(svc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{svc: svc}
}

func (h *CategoryHandler) List(c *gin.Context) {
	categories, err := h.svc.GetAll(c.Request.Context())
	if err != nil {
		errorResponse(c, http.StatusInternalServerError, "failed to get categories")
		return
	}

	c.JSON(http.StatusOK, categories)
}

func (h *CategoryHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid category id")
		return
	}

	category, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			errorResponse(c, http.StatusNotFound, "category not found")
			return
		}
		errorResponse(c, http.StatusInternalServerError, "failed to get category")
		return
	}

	c.JSON(http.StatusOK, category)
}
