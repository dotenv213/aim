package domain

import "time"

type TransactionType struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"not null" json:"title"`
	Code      string    `gorm:"unique;not null" json:"code"`
	CreatedAt time.Time `json:"created_at"`
}

type Category struct {
	ID                uint            `gorm:"primaryKey" json:"id"`
	Title             string          `gorm:"not null" json:"title"`
	TransactionTypeID uint            `gorm:"not null" json:"transaction_type_id"`
	TransactionType   TransactionType `gorm:"foreignKey:TransactionTypeID" json:"transaction_type"`
	CreatedAt         time.Time       `json:"created_at"`
}

type Contact struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	Phone     string    `json:"phone"`
	CreatedAt time.Time `json:"created_at"`
}

type Transaction struct {
	ID          uint    `gorm:"primaryKey" json:"id"`
	UserID      uint    `gorm:"index;not null" json:"user_id"`
	BankID      uint    `gorm:"index;not null" json:"bank_id"`
	Amount      float64 `gorm:"not null" json:"amount"`
	Description string  `json:"description"`

	TypeID uint            `gorm:"not null" json:"type_id"`
	Type   TransactionType `gorm:"foreignKey:TypeID" json:"type"`

	CategoryID uint     `gorm:"not null" json:"category_id"`
	Category   Category `gorm:"foreignKey:CategoryID" json:"category"`

	ContactID *uint    `json:"contact_id"`
	Contact   *Contact `gorm:"foreignKey:ContactID" json:"contact"`

	CreatedAt time.Time `json:"created_at"`
}

type TransactionRepository interface {
	Create(trx *Transaction) error
	GetByUserID(userID uint) ([]Transaction, error)
	GetTypeByCode(code string) (*TransactionType, error)
	GetCategoryByID(id uint) (*Category, error)
	CreateContact(contact *Contact) error
	GetContacts(userID uint) ([]Contact, error)
}

type TransactionService interface {
	CreateTransaction(userID uint, bankID uint, amount float64, typeCode string, categoryID uint, desc string, contactID *uint) (*Transaction, error)
	GetUserTransactions(userID uint) ([]Transaction, error)
	CreateContact(userID uint, name, phone string) (*Contact, error)
	GetContacts(userID uint) ([]Contact, error)
}
