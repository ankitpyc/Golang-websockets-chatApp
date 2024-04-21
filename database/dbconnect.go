package databases

import (
	database_model "TCPServer/database/models"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"sync"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func generateRandomChatID() (string, error) {
	// Create a byte slice to hold the random bytes
	randomBytes := make([]byte, 8)

	// Read random bytes from the crypto/rand package
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}

	// Convert the random bytes to a hexadecimal string
	randomString := hex.EncodeToString(randomBytes)

	// Return the hexadecimal string as the chat ID
	return randomString, nil
}

func ConnectToDB(wg *sync.WaitGroup) database_model.DBServer {
	fmt.Println("Connecting to Database")
	dsn := "host=localhost user=postgres password=Bablu@12345 dbname=userinfo port=5432 sslmode=disable TimeZone=Asia/Kolkata"
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer wg.Done()
	dbConn := database_model.DBServer{DB: conn}

	conn.AutoMigrate(&database_model.User{})
	conn.AutoMigrate(&database_model.Chats{})
	conn.AutoMigrate(&database_model.Message{})

	return dbConn

}
