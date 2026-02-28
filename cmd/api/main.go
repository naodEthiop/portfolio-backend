package main

import (
	"context"
	"log"
	"mime"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"portfolio-backend/internal/config"
	"portfolio-backend/internal/infrastructure/database"
	gh "portfolio-backend/internal/infrastructure/github"
	gormrepo "portfolio-backend/internal/infrastructure/repository/gorm"
	"portfolio-backend/internal/infrastructure/storage"
	"portfolio-backend/internal/transport/http/handler"
	httprouter "portfolio-backend/internal/transport/http/router"
	"portfolio-backend/internal/usecase"
	"portfolio-backend/pkg/auth"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load("configs/.env")

	_ = mime.AddExtensionType(".avif", "image/avif")

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	db, err := database.NewPostgres(cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}

	jwtManager := auth.NewJWTManager(cfg.Auth.JWTSecret, cfg.Auth.JWTIssuer, cfg.JWTExpiryDuration())
	var storageService storage.ImageStorage
	switch strings.ToLower(cfg.Storage.Provider) {
	case "supabase":
		supabaseStorage, err := storage.NewSupabaseStorage(
			cfg.Storage.SupabaseURL,
			cfg.Storage.SupabaseKey,
			cfg.Storage.SupabaseBucket,
			cfg.Storage.PublicBaseURL,
			cfg.Storage.UploadMaxBytes,
		)
		if err != nil {
			log.Fatalf("failed to init storage: %v", err)
		}
		storageService = supabaseStorage
	default:
		storageService = storage.NewLocalStorage(cfg.Storage.UploadBaseDir, cfg.Storage.UploadMaxBytes)
	}

	userRepo := gormrepo.NewUserRepository(db)
	projectRepo := gormrepo.NewProjectRepository(db)
	certificateRepo := gormrepo.NewCertificateRepository(db)
	profileRepo := gormrepo.NewProfileRepository(db)
	skillRepo := gormrepo.NewSkillRepository(db)
	socialLinkRepo := gormrepo.NewSocialLinkRepository(db)
	teachingRepo := gormrepo.NewTeachingRepository(db)
	awardRepo := gormrepo.NewAwardRepository(db)

	authService := usecase.NewAuthService(userRepo, jwtManager)
	projectService := usecase.NewProjectService(projectRepo, storageService)
	certificateService := usecase.NewCertificateService(certificateRepo, storageService)
	profileService := usecase.NewProfileService(profileRepo)
	skillService := usecase.NewSkillService(skillRepo)
	socialLinkService := usecase.NewSocialLinkService(socialLinkRepo)
	teachingService := usecase.NewTeachingService(teachingRepo)
	awardService := usecase.NewAwardService(awardRepo)
	// GitHub integration
	githubToken := cfg.Github.Token
	githubClient := gh.NewClient(githubToken, 15*time.Second)
	githubService := usecase.NewGithubService(githubClient, time.Duration(cfg.Github.CacheTTLSeconds)*time.Second, cfg.Github.PerPage)
	publicProjectsService := usecase.NewPublicProjectsService(githubService, projectRepo)

	if err := authService.BootstrapAdmin(context.Background(), cfg.Bootstrap.AdminEmail, cfg.Bootstrap.AdminPassword); err != nil {
		log.Fatalf("failed to bootstrap admin: %v", err)
	}

	h := handler.New(handler.Dependencies{
		AuthService:        authService,
		ProjectService:     projectService,
		CertificateService: certificateService,
		ProfileService:     profileService,
		SkillService:       skillService,
		SocialLinkService:  socialLinkService,
		TeachingService:    teachingService,
		AwardService:       awardService,
		GithubService:      githubService,
		PublicProjects:     publicProjectsService,
		GithubUser:         cfg.Github.User,
		GithubTopic:        cfg.Github.Topic,
	})

	r := httprouter.New(cfg, h, jwtManager)

	srv := &http.Server{
		Addr:         ":" + cfg.App.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(cfg.App.ReadTimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.App.WriteTimeoutSeconds) * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("api listening on :%s", cfg.App.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("graceful shutdown failed: %v", err)
	}
}
