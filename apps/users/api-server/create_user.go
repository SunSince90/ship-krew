package main

import (
	"encoding/json"
	"net/http"

	usersapi "github.com/SunSince90/ship-krew/users/api"
	fiber "github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
)

func createUserHandler(c *fiber.Ctx) error {
	var userToCreate usersapi.User

	if err := json.Unmarshal(c.Body(), &userToCreate); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if userToCreate.Name == "" {
		return c.SendStatus(http.StatusBadRequest)
	}

	if len(userToCreate.Name) > 25 {
		return c.SendStatus(http.StatusBadRequest)
	}

	if len(userToCreate.Bio) > 200 {
		return c.SendStatus(http.StatusBadRequest)
	}

	for _, user := range usersList {
		if userToCreate.Name == user.Name {
			return c.SendStatus(http.StatusConflict)
		}
	}

	userToCreate.ID = uuid.NewV4().String()
	usersList = append(usersList, userToCreate)

	return c.JSON(userToCreate)
}
