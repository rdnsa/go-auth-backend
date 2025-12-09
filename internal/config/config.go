package config

import (
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort         string
	JWTSecret       string
	JWTExpiredHours int
	MongoURI        MongoConfig
}

type MongoConfig struct {
	URI            string
	Database       string
	UserCollection string
}

var (
	instance *Config
	once     sync.Once
)

func Get() *Config {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using system env")
		}

		jwtHours, _ := strconv.Atoi(os.Getenv("JWT_EXPIRED_HOURS"))

		instance = &Config{
			AppPort:         getEnv("APP_PORT", "8080"),
			JWTSecret:       string(getEnv("JWT_SECRET", "supersecret")),
			JWTExpiredHours: jwtHours,
			Mongo: MongoConfig{
				URI:            getEnv("MONGO_URI", "mongodb://localhost:27017"),
				Database:       getEnv("MONGO_DB", "auth_db"),
				UserCollection: getEnv("MONGO_USER_COLLECTION", "users"),
			},
		}
	})
	return instance
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
