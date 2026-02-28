package handler

import "portfolio-backend/internal/usecase"

type Dependencies struct {
	AuthService        *usecase.AuthService
	ProjectService     *usecase.ProjectService
	CertificateService *usecase.CertificateService
	ProfileService     *usecase.ProfileService
	SkillService       *usecase.SkillService
	SocialLinkService  *usecase.SocialLinkService
	TeachingService    *usecase.TeachingService
	AwardService       *usecase.AwardService
	GithubService      *usecase.GithubService
	PublicProjects     *usecase.PublicProjectsService
	GithubUser         string
	GithubTopic        string
}

type Handler struct {
	deps Dependencies
}

func New(deps Dependencies) *Handler {
	return &Handler{deps: deps}
}
