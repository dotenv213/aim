package repository

import (
	"github.com/dotenv213/aim/account-service/internal/domain"
	"gorm.io/gorm"
)

type bankRepo struct {
	db *gorm.DB
}

func NewBankRepository(db *gorm.DB) domain.BankRepository {
	return &bankRepo{db: db}
}

func (r *bankRepo) Create(bank *domain.Bank) error {
	return r.db.Create(bank).Error
}

func (r *bankRepo) GetAllByUserID(userID uint) ([]domain.Bank, error) {
	var banks []domain.Bank
	err := r.db.Where("user_id = ?", userID).Find(&banks).Error
	return banks, err
}

func (r *bankRepo) GetByID(id uint) (*domain.Bank, error) {
	var bank domain.Bank
	err := r.db.First(&bank, id).Error
	return &bank, err
}