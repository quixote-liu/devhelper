package main

import (
	"log"

	"devhelper/internal/api"
	"devhelper/internal/config"
	"devhelper/internal/database"
	"devhelper/internal/repository"
	"devhelper/internal/service"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	if err := database.Init(cfg.DBPath); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}
	log.Println("database initialized")

	// Set initial admin if configured
	if cfg.AdminInitEmail != "" {
		userRepo := repository.NewUserRepo(database.DB)
		if err := userRepo.SetAdminByEmail(cfg.AdminInitEmail); err != nil {
			log.Printf("warning: could not set admin for %s: %v", cfg.AdminInitEmail, err)
		} else {
			log.Printf("admin role set for %s", cfg.AdminInitEmail)
		}
	}

	api.SetJWTSecret(cfg.JWTSecret)

	userRepo := repository.NewUserRepo(database.DB)
	historyRepo := repository.NewHistoryRepo(database.DB)
	schemaRepo := repository.NewSchemaRepo(database.DB)

	authSvc := service.NewAuthService(userRepo, cfg)
	jsonSvc := service.NewJsonService()

	r := gin.Default()
	api.SetupRoutes(r, cfg, authSvc, jsonSvc, userRepo, historyRepo, schemaRepo)

	if cfg.ServeStatic {
		log.Printf("serving static files from: %s", cfg.StaticFilesPath)
	}

	log.Printf("server starting on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
