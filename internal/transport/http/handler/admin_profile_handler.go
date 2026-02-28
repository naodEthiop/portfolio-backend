package handler

import (
	"net/http"

	"portfolio-backend/internal/usecase"
	"portfolio-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) UpsertProfile(c *gin.Context) {
	var input usecase.UpsertProfileInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	profile, err := h.deps.ProfileService.Upsert(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, profile)
}
