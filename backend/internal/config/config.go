package config

import (
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	DBPath            string
	JWTSecret         string
	JWTAccessExpiry   time.Duration
	JWTRefreshExpiry  time.Duration
	ServerPort        string
	CORSOrigins       string
	AdminInitEmail    string
	ServeStatic       bool
	StaticFilesPath   string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	accessExpiry, err := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	if err != nil {
		accessExpiry = 15 * time.Minute
	}
	refreshExpiry, err := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))
	if err != nil {
		refreshExpiry = 7 * 24 * time.Hour
	}

	return &Config{
		DBPath:           getEnv("DB_PATH", "./devhelper.db"),
		JWTSecret:        getEnv("JWT_SECRET", "change-this-secret"),
		JWTAccessExpiry:  accessExpiry,
		JWTRefreshExpiry: refreshExpiry,
		ServerPort:       getEnv("SERVER_PORT", "8080"),
		CORSOrigins:      getEnv("CORS_ORIGINS", "http://localhost:5173"),
		AdminInitEmail:   getEnv("ADMIN_INIT_EMAIL", ""),
		ServeStatic:      getEnv("SERVE_STATIC", "false") == "true",
		StaticFilesPath:  getEnv("STATIC_FILES_PATH", "../frontend/dist"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
