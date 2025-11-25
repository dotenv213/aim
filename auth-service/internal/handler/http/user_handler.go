package http

import (
	"net/http"

	"github.com/dotenv213/aim/auth-service/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type UserHandler struct {
	AuthService domain.AuthService 
}

func NewUserHandler(authService domain.AuthService) *UserHandler {
	return &UserHandler{
		AuthService: authService,
	}
}

func (h *UserHandler) RegisterHandler(c *fiber.Ctx) error {
	var req RegisterRequest
	
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	user, err := h.AuthService.Register(req.Email, req.Password)
	if err != nil {
		if err.Error() == "user already exists" {
			return c.Status(http.StatusConflict).JSON(fiber.Map{
				"message": "کاربری با این ایمیل از قبل وجود دارد.",
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"message": "ثبت نام با خطا مواجه شد.",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "ثبت نام با موفقیت انجام شد.",
		"user_id": user.ID,
		"email": user.Email,
	})
}

func (h *UserHandler) LoginHandler(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request body",
		})
	}

	token, err := h.AuthService.Login(req.Email, req.Password)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"message": "ایمیل یا رمز عبور نامعتبر است.",
		})
	}

	return c.JSON(fiber.Map{
		"message": "ورود موفقیت آمیز بود.",
		"token": token,
	})
}