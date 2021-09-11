package main

import (
	"strconv"
	"time"

	usersapi "github.com/SunSince90/ship-krew/users/api"
	fiber "github.com/gofiber/fiber/v2"
)

// TODO: this is just for testing and will be removed.
var usersList []usersapi.User

func init() {
	usersList = []usersapi.User{
		{
			ID:        "77cc09e9-a4be-4999-9748-a8f13faf491f",
			Name:      "uniquent",
			Bio:       "It is just me",
			CreatedAt: time.Now().AddDate(-1, 0, 0),
			UpdatedAt: time.Now().AddDate(-1, 0, 0),
		},
		{
			ID:        "7d33f8ed-2afe-4ff7-bfc6-f112c30c6e71",
			Name:      "nazoodl",
			Bio:       "This is my bio",
			CreatedAt: time.Now().AddDate(0, -10, -3),
			UpdatedAt: time.Now().AddDate(0, -2, 0),
		},
		{
			ID:        "4cc014a7-21ff-4e76-9eaf-696af8bf9a78",
			Name:      "minycat",
			Bio:       "I don't know what to say",
			CreatedAt: time.Now().AddDate(0, -1, -3),
			UpdatedAt: time.Now().AddDate(0, 0, -10),
		},
		{
			ID:        "751410dc-00c7-4b95-843b-df0ec45685ed",
			Name:      "lawli",
			Bio:       "I like turtles",
			CreatedAt: time.Now().AddDate(0, -5, 0),
			UpdatedAt: time.Now().AddDate(0, -4, 0),
		},
		{
			ID:        "27408e7e-2aa2-4542-bef7-72a9a34dab4e",
			Name:      "sunes",
			Bio:       "Ain't nobody got time for this",
			CreatedAt: time.Now().AddDate(0, -11, -4),
			UpdatedAt: time.Now().AddDate(0, -1, 8),
		},
	}
}

func listUsersHandler(c *fiber.Ctx) error {
	opts := parseListUsersOptions(
		c.Params("page", "0"),
		c.Params("per-page", "20"),
	)

	users, err := getUsersList(opts)
	if err != nil {
		_ = err
	}
	return c.JSON(users)
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
	return usersList, nil
}
