package user

import (
	"github.com/gofiber/fiber/v2"
	"message_broker/database"
	"strconv"
)

func GetUser(c *fiber.Ctx) error {
	id := c.Params("id")
	res := database.GetUserById(id)
	if len(res) == len([]byte("")) {
		return fiber.ErrBadRequest
	}
	return c.Send(res)
}

func GetAllUsers(c *fiber.Ctx) error {
	res := database.GetAllUsers()
	return c.Send(res)
}

func GetAmount(c *fiber.Ctx) error {
	id := c.Params("id")
	res := database.GetAmountById(id)
	if len(res) == len([]byte("")) {
		return fiber.ErrBadRequest
	}
	return c.Send(res)
}

func AddUser(c *fiber.Ctx) error {
	user := struct {
		Name    string `json:"username"`
		Balance int32  `json:"balance"`
	}{}
	if err := c.BodyParser(&user); err != nil {
		return fiber.ErrBadRequest
	}
	id := database.CreateUser(user.Name, user.Balance)
	if id == "" {
		return fiber.ErrBadRequest
	}
	newUser := database.User{Id: id, Name: user.Name, Balance: user.Balance}
	return c.JSON(newUser)
}

func AlterAmount(c *fiber.Ctx) error {
	userId, _ := strconv.Atoi(c.Params("id"))
	task := struct {
		Type   string `json:"type"`
		Amount int32  `json:"amount"`
	}{}
	if err := c.BodyParser(&task); err != nil {
		return fiber.ErrBadRequest
	}
	database.ChangeBalance(int64(userId), task.Type, task.Amount)

	return c.JSON("ok")
}
