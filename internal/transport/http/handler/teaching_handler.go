package handler

import (
	"net/http"

	"portfolio-backend/internal/usecase"
	"portfolio-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ListTeaching(c *gin.Context) {
	items, err := h.deps.TeachingService.ListVisible(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, items)
}

func (h *Handler) ListAdminTeaching(c *gin.Context) {
	items, err := h.deps.TeachingService.ListAdmin(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, items)
}

func (h *Handler) CreateTeaching(c *gin.Context) {
	var input usecase.CreateTeachingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	item, err := h.deps.TeachingService.Create(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, item)
}

func (h *Handler) UpdateTeaching(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	var input usecase.UpdateTeachingInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	item, err := h.deps.TeachingService.Update(c.Request.Context(), id, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, item)
}

func (h *Handler) DeleteTeaching(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.deps.TeachingService.Delete(c.Request.Context(), id); err != nil {
		handleServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
