package servers

import (
	databases "TCPServer/internal/database"
	models "TCPServer/internal/database/models"
	dto "TCPServer/internal/domain/dto"
	"TCPServer/internal/server/auth"
	"fmt"
	"sort"
	"time"

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
		var chat models.Chat
		if err := c.BodyParser(&chat); err != nil {
			return err
		}
		chats, err := db.ChatRepo.FetchChats(chat.ID)
		if err != nil {
			return err
		}
		return c.JSON(chats)
	}
}

func FetchUserChats(db *databases.DBServer) fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("Loading user chats")
		var user models.User
		if err := c.BodyParser(&user); err != nil {
			return err
		}
		chats, err := db.ChatRepo.LoadAllUserChats(user.ID)
		if err != nil {
			return err
		}
		return c.JSON(formatUserResponse(chats, user))
	}
}

func formatUserResponse(chats []models.Chat, user models.User) dto.UserDataResponse {
	buildChats(chats)
	return dto.UserDataResponse{
		UserId:       user.ID,
		UserName:     user.Username,
		ChatsHistory: chats,
	}
}

type ByLastMessageTimestamp []models.Chat

func (a ByLastMessageTimestamp) Len() int      { return len(a) }
func (a ByLastMessageTimestamp) Swap(i, j int) { a[i], a[j] = a[j], a[i] }
func (a ByLastMessageTimestamp) Less(i, j int) bool {
	// Get the timestamps of the last messages
	lastTimestampI := getLastMessageTimestamp(a[i])
	lastTimestampJ := getLastMessageTimestamp(a[j])

	// Compare the timestamps
	return lastTimestampI.After(lastTimestampJ)
}

func getLastMessageTimestamp(chat models.Chat) time.Time {
	if len(chat.Messages) == 0 {
		return time.Time{}
	}
	return chat.Messages[len(chat.Messages)-1].CreatedAt
}

func buildChats(chats []models.Chat) {
	sort.Sort(ByLastMessageTimestamp(chats))
}
