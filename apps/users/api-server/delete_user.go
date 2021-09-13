package main

import (
	"github.com/SunSince90/ship-krew/users/api"
	fiber "github.com/gofiber/fiber/v2"
)

func deleteUserHandler(c *fiber.Ctx) error {
	userName := c.Params("name")

	exists, err := getUser(api.GetUserOptions{Name: userName})
	if err != nil {
		c.SendStatus(fiber.ErrInternalServerError.Code)
		return c.Send([]byte(err.Error()))
	}

	if exists == nil {
		return c.SendStatus(fiber.ErrNotFound.Code)
	}

	if len(usersList) > 0 {
		deleteTestUser(userName)
		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendStatus(fiber.ErrNotImplemented.Code)
}
