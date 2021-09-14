package main

import (
	"encoding/json"

	"github.com/SunSince90/ship-krew/users/api-server/pkg/api"
	fiber "github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
)

func createUserHandler(c *fiber.Ctx) error {
	var userToCreate api.User

	if err := json.Unmarshal(c.Body(), &userToCreate); err != nil {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	if userToCreate.Name == "" {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	if len(userToCreate.Name) > 25 {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	if len(userToCreate.Bio) > 200 {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	for _, user := range usersList {
		if userToCreate.Name == user.Name {
			return c.SendStatus(fiber.ErrConflict.Code)
		}
	}

	userToCreate.ID = uuid.NewV4().String()
	if len(usersList) > 0 {
		createTestUsers(1)
		return c.SendStatus(fiber.StatusOK)
	}

	return c.SendStatus(fiber.ErrNotImplemented.Code)
}
