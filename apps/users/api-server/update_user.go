package main

import (
	"encoding/json"
	"net/http"

	"github.com/SunSince90/ship-krew/apps/users/api-server/pkg/api"
	fiber "github.com/gofiber/fiber/v2"
)

func updateUserHandler(c *fiber.Ctx) error {
	var newData api.User
	userName := c.Params("name")
	if err := json.Unmarshal(c.Body(), &newData); err != nil {
		return c.SendStatus(http.StatusBadRequest)
	}

	if newData.Name == "" {
		return c.SendStatus(http.StatusBadRequest)
	}

	if len(newData.Name) > 25 {
		return c.SendStatus(http.StatusBadRequest)
	}

	previous, err := getUser(api.GetUserOptions{Name: userName})
	if err != nil {
		c.SendStatus(fiber.ErrInternalServerError.Code)
		return c.Send([]byte(err.Error()))
	}

	if previous == nil {
		return c.SendStatus(fiber.ErrNotFound.Code)
	}

	newUser, err := getUser(api.GetUserOptions{Name: newData.Name})
	if err != nil {
		c.SendStatus(fiber.ErrInternalServerError.Code)
		return c.Send([]byte(err.Error()))
	}

	if newUser != nil {
		return c.SendStatus(fiber.ErrConflict.Code)
	}

	if len(newData.Bio) > 200 {
		return c.SendStatus(fiber.ErrBadRequest.Code)
	}

	if len(usersList) > 0 {
		deleteTestUser(userName)
		updateTestUser(newData.Name, newData)
	}

	return c.JSON(newData)
}
