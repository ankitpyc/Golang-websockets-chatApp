package servers

import (
	dbhandler "TCPServer/internal/database/handlers"
	models "TCPServer/internal/database/models"
	"TCPServer/internal/server/auth"
	"github.com/gofiber/fiber/v2"
)

func HandleFetchData(db *models.DBServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var userPayload models.User
		if err := c.BodyParser(&userPayload); err != nil {
			return err
		}
		user := dbhandler.CreateUser(db, &userPayload)
		return c.JSON(user)
	}
}

func LoginHandler(db *models.DBServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var userPayload models.User
		if err := c.BodyParser(&userPayload); err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
		user := dbhandler.LoginDetails(db, &userPayload)
		token := auth.GetToken("token_valid")
		if user == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		c.Set("Authorization", "Bearer "+token)
		return c.Status(200).JSON(fiber.Map{"user": user})
	}
}

func CreateUserHandler(db *models.DBServer) fiber.Handler {

	return func(c *fiber.Ctx) error {
		var payload CreateAccountRequest
		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		return c.JSON(payload)
	}
}
