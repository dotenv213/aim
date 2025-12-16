package service

import (
	"fmt"
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

func (s *bankService) GetBankByID(id uint) (*domain.Bank, error) {
	return s.repo.GetByID(id)
}

func (s *bankService) UpdateBalance(id uint, amount float64, trxType string) error {
	bank, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}

	if trxType == "deposit" {
		bank.Balance += amount
	} else if trxType == "withdraw" {
		fmt.Printf("Reducing balance for bank %d. Current: %f, Amount: %f\n", id, bank.Balance, amount)
		if bank.Balance < amount {
			return fmt.Errorf("insufficient funds")
		}
		bank.Balance -= amount
	}

	return s.repo.Update(bank)
}
