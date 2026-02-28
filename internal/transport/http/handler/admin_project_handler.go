package handler

import (
	"net/http"

	"portfolio-backend/internal/usecase"
	"portfolio-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type featuredPayload struct {
	Featured bool `json:"featured"`
}

func (h *Handler) ListAdminProjects(c *gin.Context) {
	projects, err := h.deps.ProjectService.ListAdmin(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, projects)
}

func (h *Handler) CreateProject(c *gin.Context) {
	var input usecase.CreateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	project, err := h.deps.ProjectService.Create(c.Request.Context(), input)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusCreated, project)
}

func (h *Handler) UpdateProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	var input usecase.UpdateProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	project, err := h.deps.ProjectService.Update(c.Request.Context(), id, input)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, project)
}

func (h *Handler) DeleteProject(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	if err := h.deps.ProjectService.Delete(c.Request.Context(), id); err != nil {
		handleServiceError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func (h *Handler) ToggleProjectFeatured(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	var payload featuredPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		response.Error(c, http.StatusBadRequest, "invalid request payload")
		return
	}
	project, err := h.deps.ProjectService.SetFeatured(c.Request.Context(), id, payload.Featured)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, project)
}

func (h *Handler) UploadProjectImage(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	fileHeader, err := c.FormFile("image")
	if err != nil {
		response.Error(c, http.StatusBadRequest, "missing image file")
		return
	}
	project, err := h.deps.ProjectService.UploadImage(c.Request.Context(), id, fileHeader)
	if err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, project)
}
