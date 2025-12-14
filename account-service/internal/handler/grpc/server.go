package grpc

import (
	"context"
	"fmt"

	"github.com/dotenv213/aim/account-service/internal/domain"
	pb "github.com/dotenv213/aim/account-service/proto/bank"
)

type BankGrpcHandler struct {
	pb.UnimplementedBankServiceServer
	service domain.BankService
}

func NewBankGrpcHandler(service domain.BankService) *BankGrpcHandler {
	return &BankGrpcHandler{
		service: service,
	}
}

func (h *BankGrpcHandler) GetBankAccount(ctx context.Context, req *pb.GetBankAccountRequest) (*pb.GetBankAccountResponse, error) {
	targetBank, err := h.service.GetBankByID(uint(req.BankId))
	if err != nil {
		return nil, fmt.Errorf("bank not found or db error: %v", err)
	}

	if targetBank.UserID != uint(req.UserId) {
		return nil, fmt.Errorf("access denied: bank account does not belong to this user")
	}

	return &pb.GetBankAccountResponse{
		Id:        uint64(targetBank.ID),
		Balance:   targetBank.Balance,
		OwnerName: targetBank.Name,
	}, nil
}