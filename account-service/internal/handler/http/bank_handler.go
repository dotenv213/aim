package http

import (
	"net/http"

	"github.com/dotenv213/aim/account-service/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type CreateBankRequest struct {
	Name           string  `json:"name"`
	InitialBalance float64 `json:"initial_balance"`
}

type BankHandler struct {
	Service domain.BankService
}

func NewBankHandler(service domain.BankService) *BankHandler {
	return &BankHandler{Service: service}
}

func (h *BankHandler) CreateBankHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req CreateBankRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"message": "Invalid body"})
	}

	bank, err := h.Service.CreateBank(userID, req.Name, req.InitialBalance)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to create bank"})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "Bank account created",
		"bank":    bank,
	})
}

func (h *BankHandler) GetBanksHandler(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	banks, err := h.Service.GetUserBanks(userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"message": "Failed to fetch banks"})
	}

	return c.JSON(fiber.Map{
		"banks": banks,
	})
}
