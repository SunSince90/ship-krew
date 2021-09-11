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

	user, err := getUser(opts)
	if err != nil {
		_ = err
	}

	if user == nil {
		return c.SendStatus(404)
	}

	return c.JSON(user)
}

func getUser(opts usersapi.GetUserOptions) (*usersapi.User, error) {
	for _, user := range usersList {
		if user.Name == opts.Name {
			return &user, nil
		}
	}

	return nil, nil
}
