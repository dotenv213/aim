package domain

import "time"

type Bank struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index;not null" json:"user_id"`
	Name      string    `gorm:"not null" json:"name"`
	Balance   float64   `gorm:"default:0" json:"balance"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BankRepository interface {
	Create(bank *Bank) error
	Update(bank *Bank) error
	GetAllByUserID(userID uint) ([]Bank, error)
	GetByID(id uint) (*Bank, error)
}

type BankService interface {
	CreateBank(userID uint, name string, initialBalance float64) (*Bank, error)
	GetUserBanks(userID uint) ([]Bank, error)
	GetBankByID(id uint) (*Bank, error)
	UpdateBalance(id uint, amount float64, trxType string) error
}
