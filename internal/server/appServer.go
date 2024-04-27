package servers

import (
	dbhandler "TCPServer/internal/database/handlers"
	models "TCPServer/internal/database/models"

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
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH,OPTIONS",
	}))
	defer wg.Done()

	app.Post("/fetchUserAccount", func(c *fiber.Ctx) error {
		var userPayload models.User
		if err := c.BodyParser(&userPayload); err != nil {
			return err
		}
		user := dbhandler.CreateUser(db, &userPayload)
		return c.JSON(user)
	})

	app.Post("/LoginUser", func(c *fiber.Ctx) error {
		var userPayload models.User
		if err := c.BodyParser(&userPayload); err != nil {
			return err
		}
		user := dbhandler.LoginDetails(db, &userPayload)
		return c.JSON(user)
	})

	app.Post("/createUserAccount", func(c *fiber.Ctx) error {
		var payload CreateAccountRequest
		if err := c.BodyParser(&payload); err != nil {
			return err
		}
		return c.JSON(payload)
	})

	log.Fatal(app.Listen(":3023"))
	log.Print("Listeing for http request")
}
