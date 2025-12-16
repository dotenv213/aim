package main

import (
	"log"
	"time"

	"github.com/dotenv213/aim/transaction-service/internal/domain"
	"github.com/dotenv213/aim/transaction-service/internal/handler/http"
	"github.com/dotenv213/aim/transaction-service/internal/middleware"
	"github.com/dotenv213/aim/transaction-service/internal/repository"
	"github.com/dotenv213/aim/transaction-service/internal/service"
	grpcClient "github.com/dotenv213/aim/transaction-service/pkg/client/grpc"
	"github.com/dotenv213/aim/transaction-service/pkg/config"
	"github.com/dotenv213/aim/transaction-service/pkg/database"
	"github.com/dotenv213/aim/transaction-service/pkg/rabbitmq"
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

	// --- RabbitMQ Connection ---
	rmqURL := "amqp://user:password@aim_rabbitmq:5672/"
	var producer *rabbitmq.RabbitMQProducer
	var rmqErr error

	log.Println("Connecting to RabbitMQ...")
	for i := 0; i < 30; i++ {
		producer, rmqErr = rabbitmq.NewRabbitMQProducer(rmqURL)
		if rmqErr == nil {
			log.Println("Successfully connected to RabbitMQ!")
			break
		}
		log.Printf("RabbitMQ not ready yet... retrying in 2s (Attempt %d/30)", i+1)
		time.Sleep(2 * time.Second)
	}

	if rmqErr != nil {
		log.Fatalf("Could not connect to RabbitMQ after 60 seconds: %v", rmqErr)
	}
	defer producer.Close()

	accountClient := grpcClient.NewAccountClient("aim_account_service:50051")

	trxRepo := repository.NewTransactionRepository(db)

	trxService := service.NewTransactionService(trxRepo, accountClient, producer)

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
