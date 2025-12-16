package grpc

import (
	"context"
	"fmt"
	"log"

	pb "github.com/dotenv213/aim/account-service/proto/bank"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AccountClient struct {
	client pb.BankServiceClient
}

func NewAccountClient(address string) *AccountClient {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Account Service: %v", err)
	}

	client := pb.NewBankServiceClient(conn)
	return &AccountClient{
		client: client,
	}
}

func (c *AccountClient) ValidateBankAccount(ctx context.Context, bankID uint, userID uint) (*pb.GetBankAccountResponse, error) {
	req := &pb.GetBankAccountRequest{
		BankId: uint64(bankID),
		UserId: uint64(userID),
	}

	res, err := c.client.GetBankAccount(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("error calling account service: %v", err)
	}

	return res, nil
}
