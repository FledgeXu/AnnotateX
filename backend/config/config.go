package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
}

type Config struct {
	DATABASE_URL         string
	LISTEN_ADDRESS       string
	JWT_SECRET           string
	JWT_TIMEOUT          time.Duration
	REDIS_ADDRESS        string
	REDIS_PASSWORD       string
	REDIS_DB             int
	SUPER_ADMIN_USERNAME string
	SUPER_ADMIN_PASSWORD string
	TEMP_DIR             string
	ALLOW_ORIGINS        []string
	S3Config             S3Config
}

var AppConfig *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		AppConfig = Load()
	})
	return AppConfig
}

func Load() *Config {
	c := &Config{
		DATABASE_URL:         getEnvOrPanic("DATABASE_URL"),
		LISTEN_ADDRESS:       getEnvOrDefault("LISTEN_ADDRESS", ":80"),
		JWT_SECRET:           getEnvOrDefault("JWT_SECRET", "JWT_SECRET"),
		JWT_TIMEOUT:          getDurationEnvFlexible("JWT_TIMEOUT", 72*time.Hour),
		REDIS_ADDRESS:        getEnvOrPanic("REDIS_ADDRESS"),
		REDIS_PASSWORD:       getEnvOrDefault("REDIS_PASSWORD", ""),
		REDIS_DB:             GetEnvInt("REDIS_DB", 0),
		SUPER_ADMIN_USERNAME: getEnvOrDefault("SUPER_ADMIN_USERNAME", "superadmin"),
		SUPER_ADMIN_PASSWORD: getEnvOrDefault("SUPER_ADMIN_PASSWORD", "superadmin"),
		TEMP_DIR:             getEnvOrDefault("TEMP_DIR", ""),
		ALLOW_ORIGINS: func() []string {
			parts := strings.Split(getEnvOrDefault("ALLOW_ORIGINS", ""), ",")
			for i, v := range parts {
				parts[i] = strings.TrimSpace(v)
			}
			return parts
		}(),
		S3Config: S3Config{
			Endpoint:  getEnvOrDefault("S3_ENDPOINT", "localhost:9000"),
			AccessKey: getEnvOrDefault("S3_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnvOrDefault("S3_SECRET_KEY", "minioadmin"),
			UseSSL:    getEnvOrDefaultBool("S3_USE_SSL", false),
		},
	}
	return c
}

func getEnvOrPanic(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is required but not set", key)
	}
	return value
}

func getEnvOrDefault(key, def string) string {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	return value
}

func getEnvOrDefaultBool(key string, def bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	parsed, err := strconv.ParseBool(value)
	if err != nil {
		return def
	}
	return parsed
}

func getDurationEnvFlexible(key string, def time.Duration) time.Duration {
	raw := os.Getenv(key)
	if raw == "" {
		return def
	}
	dur, err := time.ParseDuration(raw)
	if err != nil {
		log.Printf("Invalid duration for %s: %v, using default %v", key, err, def)
		return def
	}
	return dur
}

func GetEnvInt(key string, def int) int {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	if parsed, err := strconv.Atoi(value); err == nil {
		return parsed
	}
	return def
}
