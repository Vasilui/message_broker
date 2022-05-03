package main

import (
	"github.com/gofiber/fiber/v2"
	"message_broker/database"
	"message_broker/user"
)

func setupRouters(app *fiber.App) {
	app.Get("api/v1/user", user.GetAllUsers)
	app.Get("api/v1/user/:id", user.GetUser)
	app.Post("api/v1/user/", user.AddUser)
	app.Get("api/v1/user/:id/amount", user.GetAmount)
	app.Post("api/v1/user/:id/amount", user.AlterAmount)
}

func main() {
	database.InitData()
	app := fiber.New()
	setupRouters(app)
	app.Listen(":3333")

}
