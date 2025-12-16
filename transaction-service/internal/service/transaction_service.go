package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dotenv213/aim/transaction-service/internal/domain"
	grpcClient "github.com/dotenv213/aim/transaction-service/pkg/client/grpc"
	"github.com/dotenv213/aim/transaction-service/pkg/rabbitmq"
)

type transactionService struct {
	repo          domain.TransactionRepository
	accountClient *grpcClient.AccountClient
	producer      *rabbitmq.RabbitMQProducer
}

func NewTransactionService(repo domain.TransactionRepository, accClient *grpcClient.AccountClient, producer *rabbitmq.RabbitMQProducer) domain.TransactionService {
	return &transactionService{
		repo:          repo,
		accountClient: accClient,
		producer:      producer,
	}
}

func (s *transactionService) CreateTransaction(userID uint, bankID uint, amount float64, typeCode string, categoryID uint, desc string, contactID *uint) (*domain.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	bankInfo, err := s.accountClient.ValidateBankAccount(ctx, bankID, userID)
	if err != nil {
		return nil, fmt.Errorf("bank validation failed: %v", err)
	}

	trxType, err := s.repo.GetTypeByCode(typeCode)
	if err != nil {
		return nil, errors.New("invalid transaction type code")
	}

	category, err := s.repo.GetCategoryByID(categoryID)
	if err != nil {
		return nil, errors.New("invalid category id")
	}

	if category.TransactionTypeID != trxType.ID {
		return nil, errors.New("category does not match transaction type")
	}

	if trxType.Code == "withdraw" {
		if bankInfo.Balance < amount {
			return nil, errors.New("insufficient balance")
		}
	}

	trx := &domain.Transaction{
		UserID:      userID,
		BankID:      bankID,
		Amount:      amount,
		TypeID:      trxType.ID,
		CategoryID:  category.ID,
		Description: desc,
		ContactID:   contactID,
	}

	err = s.repo.Create(trx)
	if err != nil {
		return nil, err
	}

	go func() {
		pubErr := s.producer.PublishBalanceUpdate(trx.BankID, trx.Amount, trxType.Code)
		if pubErr != nil {
			fmt.Printf("FAILED to publish event: %v\n", pubErr)
		}
	}()

	return trx, nil
}

func (s *transactionService) GetUserTransactions(userID uint) ([]domain.Transaction, error) {
	return s.repo.GetByUserID(userID)
}

func (s *transactionService) CreateContact(userID uint, name, phone string) (*domain.Contact, error) {
	c := &domain.Contact{
		UserID: userID,
		Name:   name,
		Phone:  phone,
	}
	err := s.repo.CreateContact(c)
	return c, err
}

func (s *transactionService) GetContacts(userID uint) ([]domain.Contact, error) {
	return s.repo.GetContacts(userID)
}
