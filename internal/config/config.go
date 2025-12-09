package config

import (
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort         string
	JWTSecret       string
	JWTExpiredHours int

	MongoURI      string
	MongoDatabase string
	MongoUserColl string
}

var (
	instance *Config
	once     sync.Once
)

func Get() *Config {
	once.Do(func() {
		_ = godotenv.Load()

		hours, _ := strconv.Atoi(getEnv("JWT_EXPIRED_HOURS", "72"))

		instance = &Config{
			AppPort:         getEnv("APP_PORT", "8080"),
			JWTSecret:       getEnv("JWT_SECRET", "super-secret-jwt-12345678901234567890"),
			JWTExpiredHours: hours,
			MongoURI:        getEnv("MONGO_URI", "mongodb://localhost:27017"),
			MongoDatabase:   getEnv("MONGO_DB", "auth_db"),
			MongoUserColl:   getEnv("MONGO_USER_COLLECTION", "users"),
		}
	})
	return instance
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
