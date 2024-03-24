package authentication

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

var secretKey = "secret_key"

type customClaims struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	ID       uint   `json:"id"`
	jwt.StandardClaims
}

func AuthMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenString := c.Get("Authorization")[7:]
		if tokenString == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status_code": fiber.StatusUnauthorized,
				"message":     "insert token please",
			})
		}

		// verify and validate token
		claims := &customClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(secretKey), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(&fiber.Map{
				"status_code": fiber.StatusUnauthorized,
				"message":     "invalid token",
			})
		}

		c.Locals("username", claims.Username)
		c.Locals("email", claims.Email)
		c.Locals("id", claims.ID)

		return c.Next()
	}
}

func GenerateToken(username, email string, id uint) (string, error) {
	// Define token claims
	claims := jwt.MapClaims{
		"id":       id,
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	}

	// Create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate signed token
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

func GenerateHashedPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CompareHashedPassword(hashedPassword, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	return err == nil
}
