package main

import (
	"fmt"

	"github.com/SunSince90/ship-krew/apps/users/api-server/pkg/api"
	fiber "github.com/gofiber/fiber/v2"
)

func getUserHandler(c *fiber.Ctx) error {
	username := c.Params("name")
	if username == "" {
		return fmt.Errorf("username is empty")
	}

	opts := api.GetUserOptions{
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

func getUser(opts api.GetUserOptions) (*api.User, error) {
	if len(usersList) > 0 {
		return getTestUser(opts.Name), nil
	}

	return nil, fmt.Errorf("not implemented")
}

func getTestUser(name string) *api.User {
	user, exists := usersList[name]
	if !exists {
		return nil
	}

	return &user
}
