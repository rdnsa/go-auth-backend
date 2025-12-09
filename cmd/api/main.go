func main() {
	cfg := config.Get()

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(cfg.Mongo.URI))
	if err != nil {
		log.Fatal(err)
	}
	db := client.Database(cfg.Mongo.Database)

	userRepo := mongodb.NewUserRepository(db)
	userUsecase := usecase.NewUserUsecase(userRepo, cfg)

	r := gin.Default()
	handler := delivery.NewHandler(userUsecase)

	api := r.Group("/api/v1")
	{
		api.POST("/register", handler.Register)
		api.POST("/login", handler.Login)
		// api.GET("/profile", middleware.Auth(), handler.Profile) // nanti
	}

	log.Printf("Server running on port %s", cfg.AppPort)
	r.Run(":" + cfg.AppPort)
}