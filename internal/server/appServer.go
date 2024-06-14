package servers

import (
	models "TCPServer/internal/database"
	"TCPServer/internal/server/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"log"
	"sync"
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
	protected := app.Group("/api", middleware.Authorize)
	protected.Post("/api/fetchUserAccount", HandleFetchData(db))
	log.Fatal(app.Listen(":3023"))
}
