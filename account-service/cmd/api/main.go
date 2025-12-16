package main

import (
	"log"
	"net"
	"time"

	"github.com/dotenv213/aim/account-service/internal/domain"
	grpcHandler "github.com/dotenv213/aim/account-service/internal/handler/grpc"
	"github.com/dotenv213/aim/account-service/internal/handler/http"
	"github.com/dotenv213/aim/account-service/internal/middleware"
	"github.com/dotenv213/aim/account-service/internal/repository"
	"github.com/dotenv213/aim/account-service/internal/service"
	"github.com/dotenv213/aim/account-service/pkg/config"
	"github.com/dotenv213/aim/account-service/pkg/database"
	"github.com/dotenv213/aim/account-service/pkg/rabbitmq"
	pb "github.com/dotenv213/aim/account-service/proto/bank"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.LoadConfig()
	db := database.ConnectDB(cfg)
	db.AutoMigrate(&domain.Bank{})

	bankRepo := repository.NewBankRepository(db)
	bankService := service.NewBankService(bankRepo)

	go func() {
		rmqURL := "amqp://user:password@aim_rabbitmq:5672/"
		var consumer *rabbitmq.Consumer
		var err error

		for i := 0; i < 30; i++ {
			consumer, err = rabbitmq.NewConsumer(rmqURL, bankService)
			if err == nil {
				log.Println("Successfully connected to RabbitMQ!")
				break
			}
			log.Printf("RabbitMQ not ready yet... retrying in 2s (Attempt %d/30)", i+1)
			time.Sleep(2 * time.Second)
		}

		if err != nil {
			log.Fatalf("Could not connect to RabbitMQ after 60 seconds: %v", err)
			return
		}

		log.Println("RabbitMQ Consumer Connected!")
		consumer.Start()
	}()

	bankHttpHandler := http.NewBankHandler(bankService)

	bankGrpcHandler := grpcHandler.NewBankGrpcHandler(bankService)

	go func() {
		lis, err := net.Listen("tcp", ":50051")
		if err != nil {
			log.Fatalf("failed to listen: %v", err)
		}

		grpcServer := grpc.NewServer()

		pb.RegisterBankServiceServer(grpcServer, bankGrpcHandler)

		log.Println("gRPC Server listening on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	app := fiber.New()
	app.Use(logger.New())

	api := app.Group("/api/v1")
	accountGroup := api.Group("/accounts")
	accountGroup.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	accountGroup.Post("/", bankHttpHandler.CreateBankHandler)
	accountGroup.Get("/", bankHttpHandler.GetBanksHandler)

	log.Fatal(app.Listen(":8081"))
}
