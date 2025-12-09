package main

import (
	"context"
	"go-auth-backend/internal/config"
	"go-auth-backend/internal/handler/http"
	"go-auth-backend/internal/repository/mongodb"
	"go-auth-backend/internal/usecase"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg := config.Get()

	// Connect MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("Mongo connect error:", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database(cfg.MongoDatabase)

	// Dependency injection
	userRepo := mongodb.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, cfg.JWTSecret, cfg.JWTExpiredHours)
	handler := http.NewHandler(userUsecase)

	// Gin router
	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", handler.Register)
		v1.POST("/login", handler.Login)
	}

	log.Printf("Server running on :%s", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}
