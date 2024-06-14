package databases

import (
	database_model "TCPServer/internal/database/models"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func getDBConfig(db *DBServer) {
	fmt.Println("Reading DB Configuration")
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("error loading env details : ", err)
	}
	config := &database_model.DBConfig{
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_NAME:     os.Getenv("DB_NAME"),
	}
	db.Config = config
}

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

func connectToDB(server *DBServer) {
	getDBConfig(server)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Kolkata", server.Config.DB_HOST, server.Config.DB_USER, server.Config.DB_PASSWORD, server.Config.DB_NAME, server.Config.DB_PORT)
	fmt.Println("Connecting to Database")
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil || conn == nil {
		log.Printf("Unable to connect to database: %v \n", err)
		os.Exit(1)
	}
	server.DB = conn
	server.InitRepository()
}

func ConnectToDB(wg *sync.WaitGroup) DBServer {
	var server DBServer
	connectToDB(&server)
	defer wg.Done()
	log.Print("Creating tables User, Chats, Messages \n")
	connectError := server.DB.AutoMigrate(&database_model.User{}, &database_model.Chats{}, &database_model.Message{})
	if connectError != nil {
		log.Printf("Error Opening/Creating connection to databases: %v\n", connectError)
		os.Exit(1)
	}
	return server
}
