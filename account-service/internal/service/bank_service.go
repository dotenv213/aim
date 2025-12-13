package service

import (
	"github.com/dotenv213/aim/account-service/internal/domain"
)

type bankService struct {
	repo domain.BankRepository
}

func NewBankService(repo domain.BankRepository) domain.BankService {
	return &bankService{repo: repo}
}

func (s *bankService) CreateBank(userID uint, name string, initialBalance float64) (*domain.Bank, error) {
	bank := &domain.Bank{
		UserID:  userID,
		Name:    name,
		Balance: initialBalance,
	}
	
	err := s.repo.Create(bank)
	return bank, err
}

func (s *bankService) GetUserBanks(userID uint) ([]domain.Bank, error) {
	return s.repo.GetAllByUserID(userID)
}