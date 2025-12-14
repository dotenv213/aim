package http

import (
	"net/http"

	"github.com/dotenv213/aim/transaction-service/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type CreateTransactionRequest struct {
	BankID      uint    `json:"bank_id"`
	Amount      float64 `json:"amount"`
	TypeCode    string  `json:"type_code"` 
	CategoryID  uint    `json:"category_id"`
	Description string  `json:"description"`
	ContactID   *uint   `json:"contact_id"` 
}

type TransactionHandler struct {
	service domain.TransactionService
}

func NewTransactionHandler(service domain.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

func (h *TransactionHandler) CreateHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint) 

	var req CreateTransactionRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid body"})
	}

	trx, err := h.service.CreateTransaction(
		userID,
		req.BankID,
		req.Amount,
		req.TypeCode,
		req.CategoryID,
		req.Description,
		req.ContactID, 
	)

	if err != nil {
		if err.Error() == "insufficient balance" {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "موجودی کافی نیست!"})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": err.Error()})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Transaction created successfully",
		"transaction": trx,
	})
}

func (h *TransactionHandler) GetListHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	trxs, err := h.service.GetUserTransactions(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch transactions"})
	}

	return c.JSON(fiber.Map{
		"transactions": trxs,
	})
}

type CreateContactRequest struct {
    Name  string `json:"name"`
    Phone string `json:"phone"`
}

func (h *TransactionHandler) CreateContactHandler(c *fiber.Ctx) error {
    userID := c.Locals("user_id").(uint)
    var req CreateContactRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid body"})
    }
    
    contact, err := h.service.CreateContact(userID, req.Name, req.Phone)
    if err != nil {
         return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed"})
    }
    
    return c.Status(http.StatusCreated).JSON(contact)
}