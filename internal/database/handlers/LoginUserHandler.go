package databases

import (
	models "TCPServer/internal/database/models"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func LoginDetails(dbServer *models.DBServer, user *models.User) (users *models.User) {
	userLogin := &models.User{}
	result := dbServer.DB.Where("email = ?", user.Email).First(&userLogin)
	if result.Error != nil {
		log.Println("Error ", result.Error)
		return nil
	}
	if err := bcrypt.CompareHashAndPassword([]byte(userLogin.GetPassword()), []byte(user.GetPassword())); err != nil {
		log.Println("Password does not match")
		return nil
	}
	return userLogin
}
