package main

import (
	"log"

	"github.com/dotenv213/aim/transaction-service/internal/domain"
	"github.com/dotenv213/aim/transaction-service/internal/handler/http"
	"github.com/dotenv213/aim/transaction-service/internal/middleware"
	"github.com/dotenv213/aim/transaction-service/internal/repository"
	"github.com/dotenv213/aim/transaction-service/internal/service"
	grpcClient "github.com/dotenv213/aim/transaction-service/pkg/client/grpc"
	"github.com/dotenv213/aim/transaction-service/pkg/config"
	"github.com/dotenv213/aim/transaction-service/pkg/database"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	cfg := config.LoadConfig()

	db := database.ConnectDB(cfg)

	err := db.AutoMigrate(
		&domain.TransactionType{},
		&domain.Contact{},
		&domain.Category{},
		&domain.Transaction{},
	)
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}

	database.SeedDatabase(db)

	accountClient := grpcClient.NewAccountClient("aim_account_service:50051")

	trxRepo := repository.NewTransactionRepository(db)
	trxService := service.NewTransactionService(trxRepo, accountClient)
	trxHandler := http.NewTransactionHandler(trxService)

	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api/v1")

	transactionGroup := api.Group("/transactions")
	transactionGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret)) 
	transactionGroup.Post("/", trxHandler.CreateHandler)
	transactionGroup.Get("/", trxHandler.GetListHandler)

	contactGroup := api.Group("/contacts")
	contactGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	contactGroup.Post("/", trxHandler.CreateContactHandler)

	log.Fatal(app.Listen(":8080")) 
}