package middleware

import (
	"TCPServer/internal/server/auth"
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm/logger"
)

func Authorize(c *fiber.Ctx) error {
	log.Println(logger.Cyan, "Authorize request received")
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		log.Println(logger.Red, fmt.Errorf("Empty Authorization header"))
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Authorization header is empty"})
	}
	token := authHeader[7:]
	isValid, err := auth.ValidateToken("token_valid", token)
	if !isValid || err != nil {
		log.Println(logger.BlueBold, fmt.Errorf("Empty Authorization header"))

		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err})
	}
	return c.Next()
}
