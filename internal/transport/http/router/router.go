package router

import (
	"time"

	"portfolio-backend/internal/config"
	"portfolio-backend/internal/domain/entities"
	"portfolio-backend/internal/transport/http/handler"
	"portfolio-backend/internal/transport/http/middleware"
	"portfolio-backend/pkg/auth"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func New(cfg *config.Config, h *handler.Handler, jwtManager *auth.JWTManager) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	r.Static("/uploads", cfg.Storage.UploadBaseDir)

	r.GET("/health", h.Health)

	// public (unversioned) endpoints consumed by the portfolio homepage
	r.GET("/api/projects", h.ListProjects)
	r.GET("/api/certificates", h.ListCertificates)
	r.GET("/api/skills", h.ListSkills)
	r.GET("/api/contact", h.GetContact)
	r.GET("/api/teaching", h.ListTeaching)
	r.GET("/api/awards", h.ListAwards)

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/login", h.Login)
		v1.GET("/profile", h.GetProfile)
		v1.GET("/projects", h.ListPublicProjects)
		v1.GET("/projects/:id", h.GetProjectByID)
		v1.GET("/certificates", h.ListCertificates)
		v1.GET("/skills", h.ListSkills)
		v1.GET("/social-links", h.ListVisibleSocialLinks)
		v1.GET("/teaching", h.ListTeaching)
		v1.GET("/awards", h.ListAwards)

		admin := v1.Group("/admin")
		admin.Use(middleware.AuthRequired(jwtManager), middleware.RequireRole(entities.RoleAdmin))
		{
			admin.GET("/projects", h.ListAdminProjects)
			admin.POST("/projects", h.CreateProject)
			admin.PUT("/projects/:id", h.UpdateProject)
			admin.DELETE("/projects/:id", h.DeleteProject)
			admin.PATCH("/projects/:id/featured", h.ToggleProjectFeatured)
			admin.POST("/projects/:id/image", h.UploadProjectImage)

			admin.POST("/certificates", h.CreateCertificate)
			admin.PUT("/certificates/:id", h.UpdateCertificate)
			admin.DELETE("/certificates/:id", h.DeleteCertificate)
			admin.POST("/certificates/:id/image", h.UploadCertificateImage)

			admin.PUT("/profile", h.UpsertProfile)

			admin.POST("/skills", h.CreateSkill)
			admin.PUT("/skills/:id", h.UpdateSkill)
			admin.DELETE("/skills/:id", h.DeleteSkill)

			admin.GET("/social-links", h.ListAdminSocialLinks)
			admin.POST("/social-links", h.CreateSocialLink)
			admin.PUT("/social-links/:id", h.UpdateSocialLink)
			admin.DELETE("/social-links/:id", h.DeleteSocialLink)

			admin.GET("/teaching", h.ListAdminTeaching)
			admin.POST("/teaching", h.CreateTeaching)
			admin.PUT("/teaching/:id", h.UpdateTeaching)
			admin.DELETE("/teaching/:id", h.DeleteTeaching)

			admin.GET("/awards", h.ListAdminAwards)
			admin.POST("/awards", h.CreateAward)
			admin.PUT("/awards/:id", h.UpdateAward)
			admin.DELETE("/awards/:id", h.DeleteAward)
		}
	}

	return r
}
