package httpservice

import (
	"ebookstore/internal/model/request"
	"ebookstore/internal/model/response"
	"ebookstore/internal/service"
	"errors"
	"regexp"
	"unicode"

	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	customerService service.ICustomerService
}

func NewCustomerHandler(customerService service.ICustomerService) *CustomerHandler {
	return &CustomerHandler{
		customerService: customerService,
	}
}

func (h *CustomerHandler) Register(c *fiber.Ctx) error {
	customer := request.Register{}
	err := c.BodyParser(&customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Customer{
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	err = isValidUsernameEmailPassword(customer)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(response.Customer{
			StatusCode: fiber.StatusBadRequest,
			Message:    err.Error(),
		})
	}

	token, err := h.customerService.Register(c.Context(), &customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Customer{
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.Customer{
		StatusCode: fiber.StatusOK,
		Message:    "success",
		Token:      token,
	})
}

func (h *CustomerHandler) Login(c *fiber.Ctx) error {
	login := request.Login{}
	err := c.BodyParser(&login)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Customer{
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	token, err := h.customerService.Login(c.Context(), login)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(response.Customer{
			StatusCode: fiber.StatusInternalServerError,
			Message:    err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(response.Customer{
		StatusCode: fiber.StatusOK,
		Message:    "success",
		Token:      token,
	})
}

func isValidUsernameEmailPassword(customer request.Register) error {
	if !isValidUsername(customer.Username) {
		return errors.New("invalid username, between 4 and 16 characters in length contains only alphanumeric characters or underscores")
	}

	if !isValidEmail(customer.Email) {
		return errors.New("invalid email combination")
	}

	if !isValidPassword(customer.Password) {
		return errors.New("invalid password, at least 8 characters in length contains at least one lowercase letter, one uppercase letter, one digit, and one special character")
	}

	return nil
}

func isValidUsername(username string) bool {
	// Validate username format
	match, err := regexp.MatchString("^[a-zA-Z0-9_]{4,16}$", username)
	if err != nil {
		return false
	}

	return match
}

// CheckPasswordStrength checks if a password meets the following criteria:
// - At least 8 characters in length
// - Contains at least one lowercase letter, one uppercase letter, one digit, and one special character
func isValidPassword(password string) bool {
	var (
		hasLower   = false
		hasUpper   = false
		hasDigit   = false
		hasSpecial = false
	)

	if len(password) < 8 {
		return false
	}

	for _, c := range password {
		switch {
		case unicode.IsLower(c):
			hasLower = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsDigit(c):
			hasDigit = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}

	return hasLower && hasUpper && hasDigit && hasSpecial
}

func isValidEmail(email string) bool {
	// Validate email format
	match, err := regexp.MatchString(`^([a-zA-Z0-9_\-\.]+)@([a-zA-Z0-9-]+\.)+([a-zA-Z]{2,10})$`, email)
	if err != nil {
		return false
	}

	return match
}
