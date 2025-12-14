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

	log.Println("Seeding database with default types and categories...")

	income := domain.TransactionType{Title: "Income", Code: "INC"}
	expense := domain.TransactionType{Title: "Expense", Code: "EXP"}

	db.Create(&income)
	db.Create(&expense)

	categories := []domain.Category{
		{Title: "Salary", TransactionTypeID: income.ID},
		{Title: "Gift", TransactionTypeID: income.ID},
		
		{Title: "Food & Dining", TransactionTypeID: expense.ID},
		{Title: "Rent", TransactionTypeID: expense.ID},
		{Title: "Transportation", TransactionTypeID: expense.ID},
		{Title: "Shopping", TransactionTypeID: expense.ID},
	}

	db.Create(&categories)
	log.Println("Database seeding completed!")
}