package database

import (
	"fmt"
	"log"
	"time" 

	"github.com/dotenv213/aim/transaction-service/pkg/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB(cfg *config.Config) *gorm.DB {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort,
	)

	counts := 0
	for {
		db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err != nil {
			log.Println("Postgres not yet ready...")
			counts++
		} else {
			log.Println("Connected to Database successfully!")
			return db
		}

		if counts > 10 { 
			log.Println(err)
			log.Fatalf("Could not connect to the database after multiple retries")
		}

		log.Println("Backing off for two seconds...")
		time.Sleep(2 * time.Second) 
		continue
	}
}