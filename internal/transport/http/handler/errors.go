package handler

import (
	"errors"
	"net/http"

	"portfolio-backend/internal/usecase"
	"portfolio-backend/pkg/response"

	"github.com/gin-gonic/gin"
)

func handleServiceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, usecase.ErrNotFound):
		response.Error(c, http.StatusNotFound, "resource not found")
	case errors.Is(err, usecase.ErrUnauthorized):
		response.Error(c, http.StatusUnauthorized, "invalid credentials")
	default:
		response.Error(c, http.StatusInternalServerError, "internal server error")
	}
}
