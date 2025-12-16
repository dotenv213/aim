package main

import (
	"log"

	"github.com/dotenv213/aim/auth-service/internal/domain"
	"github.com/dotenv213/aim/auth-service/internal/handler/http"
	"github.com/dotenv213/aim/auth-service/internal/repository"
	"github.com/dotenv213/aim/auth-service/internal/service"
	"github.com/dotenv213/aim/auth-service/pkg/config"
	"github.com/dotenv213/aim/auth-service/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Load configs
	cfg := config.LoadConfig()

	// Connect db
	db := database.ConnectDB(cfg)

	// auto migrate
	err := db.AutoMigrate(&domain.User{})
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	//Dependency Injection
	// connect repository to db
	userRepo := repository.NewUserRepository(db)

	// service
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)

	// handler connect to service
	userHandler := http.NewUserHandler(authService)

	// routes
	app := fiber.New()

	app.Use(logger.New())

	api := app.Group("/api/v1")
	auth := api.Group("/auth")

	auth.Post("/register", userHandler.RegisterHandler)
	auth.Post("/login", userHandler.LoginHandler)

	// Run server
	log.Fatal(app.Listen(":8080"))
}
