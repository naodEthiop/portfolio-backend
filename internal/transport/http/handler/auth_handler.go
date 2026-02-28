package handler

import (
	"net/http"

	"portfolio-backend/internal/usecase"
	"portfolio-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func (h *Handler) Login(c *gin.Context) {
	var input usecase.LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload")
		return
	}

	result, err := h.deps.AuthService.Login(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, result)
}
