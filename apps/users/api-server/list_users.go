package main

import (
	"strconv"

	usersapi "github.com/SunSince90/ship-krew/users/api"
	fiber "github.com/gofiber/fiber/v2"
)

func listUsersHandler(c *fiber.Ctx) error {
	opts := parseListUsersOptions(
		c.Params("page", "0"),
		c.Params("per-page", "20"),
	)
	_ = opts

	if len(usersList) > 0 {
		return c.JSON(getTestUsersList())
	}

	return c.SendStatus(fiber.ErrNotImplemented.Code)
}

func parseListUsersOptions(pageParam, perPageParam string) usersapi.ListUsersOptions {
	l := log

	page, err := strconv.Atoi(pageParam)
	if err != nil {
		l.Err(err).Str("page", pageParam).Int("default", 0).Msg("invalid value provided for 'page', using default...")
		page = 0
	}

	perPage, err := strconv.Atoi(perPageParam)
	if err != nil {
		l.Err(err).Str("per-page", perPageParam).Int("default", 20).Msg("invalid value provided for 'per-page', using default...")
		perPage = 20
	}

	return usersapi.ListUsersOptions{
		Page:    page,
		PerPage: perPage,
	}
}

func getUsersList(opts usersapi.ListUsersOptions) ([]usersapi.User, error) {
	// Are we running in test mode?
	if len(usersList) > 0 {
		return getTestUsersList(), nil
	}

	return []usersapi.User{}, nil
}

func getTestUsersList() []usersapi.User {
	list := []usersapi.User{}

	for _, usr := range usersList {
		list = append(list, usr)
	}

	return list
}
