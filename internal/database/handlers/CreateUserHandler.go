package databases

import (
	"TCPServer/internal/database"
	models "TCPServer/internal/database/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CreateUser(dbServer *databases.DBServer, user *models.User) *models.User {
	hashedPassword, _ := hashPassword(user.GetPassword())
	userData := &models.User{Username: user.Username, Email: user.Email}
	userData.SetPassword(hashedPassword)
	result := dbServer.DB.Create(&userData)
	if result.Error != nil {
		log.Println("Error ", result.Error)
	}
	return userData
}
