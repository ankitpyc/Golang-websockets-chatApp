package servers

import (
	"TCPServer/internal/database"
	models "TCPServer/internal/database/models"
	"TCPServer/internal/server/auth"
	"github.com/gofiber/fiber/v2"
)

func LoginHandler(db *databases.DBServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var userPayload models.User
		if err := c.BodyParser(&userPayload); err != nil {
			c.Status(fiber.StatusBadRequest)
			return err
		}
		user, err := db.UserRepo.Login(&userPayload)
		if err != nil && err.Error() == "invalid User Name and Password" {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		if user == nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		//TODO : Secret Needs to be Externalized
		token := auth.GetToken("token_valid")
		c.Set("Authorization", "Bearer "+token)
		return c.Status(200).JSON(fiber.Map{"user": user})
	}
}

func CreateUserHandler(db *databases.DBServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var userPayload models.User
		if err := c.BodyParser(&userPayload); err != nil {
			return err
		}
		user, _ := db.UserRepo.CreateUser(&userPayload)
		return c.JSON(user)
	}
}

func FetchAllChats(db *databases.DBServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var messages models.Message
		if err := c.BodyParser(&messages); err != nil {
			return err
		}
		chats, err := db.ChatRepo.FetchChats(messages.ChatId)
		if err != nil {
			return err
		}
		return c.JSON(chats)
	}
}
