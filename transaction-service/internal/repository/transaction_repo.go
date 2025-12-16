package repository

import (
	"github.com/dotenv213/aim/transaction-service/internal/domain"
	"gorm.io/gorm"
)

type transactionRepo struct {
	db *gorm.DB
}

func NewTransactionRepository(db *gorm.DB) domain.TransactionRepository {
	return &transactionRepo{db: db}
}

func (r *transactionRepo) Create(trx *domain.Transaction) error {
	return r.db.Create(trx).Error
}

func (r *transactionRepo) GetByUserID(userID uint) ([]domain.Transaction, error) {
	var trxs []domain.Transaction
	err := r.db.Preload("Type").Preload("Category").Preload("Contact").
		Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&trxs).Error
	return trxs, err
}

func (r *transactionRepo) GetTypeByCode(code string) (*domain.TransactionType, error) {
	var t domain.TransactionType
	err := r.db.Where("code = ?", code).First(&t).Error
	return &t, err
}

func (r *transactionRepo) GetCategoryByID(id uint) (*domain.Category, error) {
	var c domain.Category
	err := r.db.First(&c, id).Error
	return &c, err
}

func (r *transactionRepo) CreateContact(contact *domain.Contact) error {
	return r.db.Create(contact).Error
}

func (r *transactionRepo) GetContacts(userID uint) ([]domain.Contact, error) {
	var contacts []domain.Contact
	err := r.db.Where("user_id = ?", userID).Find(&contacts).Error
	return contacts, err
}
