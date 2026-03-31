package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI           string
	MongoDB            string
	JWTSecret          string
	Port               string // default "8081"
	Framework          string // hardcode "Gin"
	TemplateDir        string // default "../shared/templates"
	IsProduction       bool   // APP_ENV == "production"
	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

func Load() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		MongoURI:           os.Getenv("MONGO_URI"),
		MongoDB:            os.Getenv("MONGO_DB"),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		Port:               getEnv("PORT", "8081"),
		Framework:          "Gin",
		TemplateDir:        getEnv("TEMPLATE_DIR", "../shared/templates"),
		IsProduction:       os.Getenv("APP_ENV") == "production",
		GoogleClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		GoogleRedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
	}

	if cfg.MongoURI == "" || cfg.MongoDB == "" || cfg.JWTSecret == "" {
		log.Fatal("Missing required env: MONGO_URI, MONGO_DB, JWT_SECRET")
	}

	return cfg
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
