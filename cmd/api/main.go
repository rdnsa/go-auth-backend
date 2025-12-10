package main

import (
	"context"
	"log"
	"time"

	"go-auth-backend/internal/config"
	"go-auth-backend/internal/handler/http"
	"go-auth-backend/internal/repository/mongodb"
	"go-auth-backend/internal/usecase"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Ambil config dari .env
	cfg := config.Get()

	// Koneksi ke MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatal("Gagal konek ke MongoDB:", err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			log.Println("Error disconnect MongoDB:", err)
		}
	}()

	// Test koneksi
	if err = client.Ping(ctx, nil); err != nil {
		log.Fatal("MongoDB tidak respons:", err)
	}
	log.Println("MongoDB connected!")

	db := client.Database(cfg.MongoDatabase)

	// Dependency Injection
	userRepo := mongodb.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, cfg.JWTSecret, cfg.JWTExpiredHours)
	handler := http.NewHandler(userUsecase)

	// === GIN SETUP + CORS (INI YANG BIKIN BISA DARI localhost:3000) ===
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// CORS untuk development (izinin frontend Next.js)
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // frontend kamu
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// API Routes
	v1 := r.Group("/api/v1")
	{
		v1.POST("/register", handler.Register)
		v1.POST("/login", handler.Login)
	}

	// Jalankan server
	log.Printf("Server Go berjalan di http://localhost:%s", cfg.AppPort)
	log.Println("Silakan buka frontend di http://localhost:3000")
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatal("Server error:", err)
	}
}
