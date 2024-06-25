package servers

import (
	models "TCPServer/internal/database"
	"log"
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func StartWebServer(wg *sync.WaitGroup, db *models.DBServer) {
	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowHeaders:     "Origin,Content-Type,Accept,Content-Length,Accept-Language,Accept-Encoding,Connection,Access-Control-Allow-Origin",
		AllowOrigins:     "http://localhost:3000",
		AllowCredentials: true,
		ExposeHeaders:    "Authorization",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	defer wg.Done()
	app.Post("/LoginUser", LoginHandler(db))
	app.Post("/api/createUserAccount", CreateUserHandler(db))
	app.Post("/api/fetchAllChats", FetchAllChats(db))
	app.Post("/api/LoadUserChats", FetchUserChats(db))

	log.Fatal(app.Listen(":3023"))
}
