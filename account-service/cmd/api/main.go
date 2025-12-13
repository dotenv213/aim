package main

import (
	"log"

	"github.com/dotenv213/aim/account-service/internal/domain"
	"github.com/dotenv213/aim/account-service/internal/handler/http"
	"github.com/dotenv213/aim/account-service/internal/middleware"
	"github.com/dotenv213/aim/account-service/internal/repository"
	"github.com/dotenv213/aim/account-service/internal/service"
	"github.com/dotenv213/aim/account-service/pkg/config"
	"github.com/dotenv213/aim/account-service/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.LoadConfig()
	db := database.ConnectDB(cfg)

	db.AutoMigrate(&domain.Bank{})

	bankRepo := repository.NewBankRepository(db)
	bankService := service.NewBankService(bankRepo)
	bankHandler := http.NewBankHandler(bankService)

	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api/v1")
	
	accountGroup := api.Group("/accounts")
	
	accountGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))

	accountGroup.Post("/", bankHandler.CreateBankHandler)
	accountGroup.Get("/", bankHandler.GetBanksHandler)

	log.Fatal(app.Listen(":8081")) 
}