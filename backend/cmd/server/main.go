package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lcs/devhelper/internal/api"
	"github.com/lcs/devhelper/internal/config"
	"github.com/lcs/devhelper/internal/database"
	"github.com/lcs/devhelper/internal/repository"
	"github.com/lcs/devhelper/internal/service"
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
	api.SetupRoutes(r, cfg.CORSOrigins, authSvc, jsonSvc, userRepo, historyRepo, schemaRepo)

	log.Printf("server starting on :%s", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
