package handler

import (
	"errors"
	"net/http"

	"portfolio-backend/internal/usecase"
	"portfolio-backend/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (h *Handler) Health(c *gin.Context) {
	response.JSON(c, http.StatusOK, gin.H{"status": "ok"})
}

func (h *Handler) GetProfile(c *gin.Context) {
	profile, err := h.deps.ProfileService.Get(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, profile)
}

func (h *Handler) ListPublicProjects(c *gin.Context) {
	projects, err := h.deps.ProjectService.ListPublic(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, projects)
}

func (h *Handler) GetProjectByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		response.Error(c, http.StatusBadRequest, "invalid id")
		return
	}
	project, err := h.deps.ProjectService.GetByID(c.Request.Context(), id)
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, project)
}

func (h *Handler) ListCertificates(c *gin.Context) {
	certificates, err := h.deps.CertificateService.List(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, certificates)
}

func (h *Handler) ListSkills(c *gin.Context) {
	skills, err := h.deps.SkillService.List(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, skills)
}

func (h *Handler) ListVisibleSocialLinks(c *gin.Context) {
	links, err := h.deps.SocialLinkService.ListVisible(c.Request.Context())
	if err != nil {
		handleServiceError(c, err)
		return
	}
	response.JSON(c, http.StatusOK, links)
}

func (h *Handler) GetContact(c *gin.Context) {
	ctx := c.Request.Context()

	profile, err := h.deps.ProfileService.Get(ctx)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	links, err := h.deps.SocialLinkService.ListVisible(ctx)
	if err != nil {
		handleServiceError(c, err)
		return
	}

	response.JSON(c, http.StatusOK, gin.H{
		"profile":      profile,
		"social_links": links,
	})
}

func (h *Handler) ListProjects(c *gin.Context) {
	ctx := c.Request.Context()
	user := h.deps.GithubUser
	topic := h.deps.GithubTopic
	if user == "" {
		user = "naodEthiop"
	}
	if topic == "" {
		topic = "portfolio"
	}
	if h.deps.PublicProjects == nil {
		response.Error(c, http.StatusInternalServerError, "projects service not configured")
		return
	}
	out, err := h.deps.PublicProjects.List(ctx, user, topic)
	if err != nil {
		switch {
		case errors.Is(err, usecase.ErrRateLimited):
			response.Error(c, http.StatusTooManyRequests, "github rate limit exceeded")
		case errors.Is(err, usecase.ErrUpstream):
			response.Error(c, http.StatusBadGateway, "failed to fetch github projects")
		default:
			response.Error(c, http.StatusInternalServerError, "internal server error")
		}
		return
	}
	response.JSON(c, http.StatusOK, out)
}
