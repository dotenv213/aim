package database

import (
	"log"

	"github.com/dotenv213/aim/transaction-service/internal/domain"
	"gorm.io/gorm"
)

func SeedDatabase(db *gorm.DB) {
	var count int64
	db.Model(&domain.TransactionType{}).Count(&count)

	if count > 0 {
		log.Println("Database already seeded.")
		return
	}

	log.Println("Seeding database with default categories...")

	deposit := domain.TransactionType{Title: "واریز", Code: "deposit"}
	withdraw := domain.TransactionType{Title: "برداشت", Code: "withdraw"}

	db.Create(&deposit)
	db.Create(&withdraw)

	categories := []domain.Category{
		{Title: "دریافت قرض / وام", TransactionTypeID: deposit.ID}, // ID: 1
		{Title: "تسویه با بدهکار", TransactionTypeID: deposit.ID},  // ID: 2
		{Title: "سایر واریزها", TransactionTypeID: deposit.ID},     // ID: 3

		{Title: "اعطای قرض / وام به شخص", TransactionTypeID: withdraw.ID}, // ID: 4
		{Title: "تسویه با بستانکار", TransactionTypeID: withdraw.ID},      // ID: 5
		{Title: "هزینه خوراک", TransactionTypeID: withdraw.ID},            // ID: 6
		{Title: "سایر برداشت‌ها", TransactionTypeID: withdraw.ID},         // ID: 7
	}

	db.Create(&categories)
	log.Println("Database seeding completed successfully!")
}
