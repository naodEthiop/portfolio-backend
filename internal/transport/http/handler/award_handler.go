package handler

import (
	"net/http"

	"portfolio-backend/internal/usecase"
	"portfolio-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) ListAwards(c *gin.Context) {
	items, err := h.deps.AwardService.ListVisible(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, items)
}

func (h *Handler) ListAdminAwards(c *gin.Context) {
	items, err := h.deps.AwardService.ListAdmin(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, items)
}

func (h *Handler) CreateAward(c *gin.Context) {
	var input usecase.CreateAwardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	item, err := h.deps.AwardService.Create(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, item)
}

func (h *Handler) UpdateAward(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	var input usecase.UpdateAwardInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	item, err := h.deps.AwardService.Update(c.Request.Context(), id, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, item)
}

func (h *Handler) DeleteAward(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.deps.AwardService.Delete(c.Request.Context(), id); err != nil {
		handleServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}
