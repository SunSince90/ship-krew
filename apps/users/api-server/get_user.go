package main

import (
	"fmt"

	usersapi "github.com/SunSince90/ship-krew/users/api"
	fiber "github.com/gofiber/fiber/v2"
)

func getUserHandler(c *fiber.Ctx) error {
	username := c.Params("name")
	if username == "" {
		return fmt.Errorf("username is empty")
	}

	opts := usersapi.GetUserOptions{
		Name: username,
	}

	usr, err := getUser(opts)
	if err != nil {
		c.SendStatus(fiber.ErrInternalServerError.Code)
		return c.Send([]byte(err.Error()))
	}

	if usr == nil {
		return c.JSON(fiber.ErrNotFound)
	}

	return c.JSON(usr)
}

func getUser(opts usersapi.GetUserOptions) (*usersapi.User, error) {
	if len(usersList) > 0 {
		return getTestUser(opts.Name), nil
	}

	return nil, fmt.Errorf("not implemented")
}

func getTestUser(name string) *usersapi.User {
	user, exists := usersList[name]
	if !exists {
		return nil
	}

	return &user
}
