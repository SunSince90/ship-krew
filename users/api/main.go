package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"

	udb "github.com/SunSince90/ship-krew/users/api/internal/database"
	"github.com/SunSince90/ship-krew/users/api/pkg/api"
	"github.com/SunSince90/ship-krew/users/api/pkg/database"
	uerrors "github.com/SunSince90/ship-krew/users/api/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

const (
	fiberAppName string = "Users API Server"
)

var (
	log zerolog.Logger
)

func main() {
	var (
		verbosity  int
		dbSettings = &database.Settings{}
	)

	flag.IntVar(&verbosity, "verbosity", 1, "the verbosity level")

	flag.StringVar(&dbSettings.Name, "database-name", "", "the name of the database to connect to")
	flag.StringVar(&dbSettings.User, "database-user", "", "the username to connect as")
	flag.StringVar(&dbSettings.Password, "database-password", "", "the password to use for the provided user")
	flag.StringVar(&dbSettings.Address, "database-address", "localhost", "the address where mysql is running")
	flag.IntVar(&dbSettings.Port, "database-port", 3306, "the port mysql is exposing")
	flag.StringVar(&dbSettings.Charset, "database-charset", "utf8mb4", "the charset used by the database")
	flag.DurationVar(&dbSettings.ReadTimeout, "database-readtimeout", 2*time.Minute, "the charset used by the database")
	flag.DurationVar(&dbSettings.WriteTimeout, "database-writetimeout", 2*time.Minute, "the charset used by the database")
	flag.Parse()

	log = zerolog.New(os.Stderr).With().Logger()
	log.Info().Int("verbosity", verbosity).Msg("starting...")

	{
		logLevels := [4]zerolog.Level{zerolog.DebugLevel, zerolog.InfoLevel, zerolog.ErrorLevel}
		log = log.Level(logLevels[verbosity])
	}

	db, err := database.NewDatabaseConnection(dbSettings)
	if err != nil {
		log.Err(err).Msg("error while establishing connection to the database")
		return
	}

	usersDB := &udb.Database{DB: db, Logger: log}

	app := fiber.New(fiber.Config{
		AppName:               fiberAppName,
		ReadTimeout:           time.Minute,
		DisableStartupMessage: verbosity > 0,
	})

	users := app.Group("/users")

	users.Get("/", func(c *fiber.Ctx) error {
		filters := &udb.ListFilters{}

		page, err := func() (int, error) {
			p, err := url.QueryUnescape(c.Query("page", "1"))
			if err != nil {
				return 0, fmt.Errorf("error while unescaping page: %w", err)
			}

			return strconv.Atoi(p)
		}()
		if err != nil || page < 1 {
			return c.Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeInvalidPage,
					Message: uerrors.MessageInvalidPage,
					Err:     err,
				})
		}
		filters.Page = &page

		nameIn, err := url.QueryUnescape(c.Query("usernameIn"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeInvalidNameIn,
					Message: uerrors.MessageInvalidNameIn,
					Err:     err,
				})
		}

		emailIn, err := url.QueryUnescape(c.Query("emailIn"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeInvalidEmailIn,
					Message: uerrors.MessageInvalidEmailIn,
					Err:     err,
				})
		}

		idIn, err := url.QueryUnescape(c.Query("idIn"))
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeInvalidIdIn,
					Message: uerrors.MessageInvalidIdIn,
					Err:     err,
				})
		}

		switch {
		case idIn != "":
			filters.IDIn = func() (filteredIds []int64) {
				ids := strings.Split(idIn, ",")
				for _, id := range ids {
					if val, err := strconv.ParseInt(id, 10, 64); err == nil && val > 0 {
						filteredIds = append(filteredIds, val)
					}
				}
				return
			}()
		case nameIn != "":
			filters.UsernameIn = strings.Split(nameIn, ",")
		case emailIn != "":
			filters.EmailIn = strings.Split(emailIn, ",")
		}

		users, err := usersDB.ListUsers(filters)
		if err != nil {
			code := err.(*uerrors.Error).Code

			return c.
				Status(uerrors.ToHTTPStatusCode(code)).
				JSON(err)
		}

		return c.JSON(users)
	})

	users.Get("/username/:username", func(c *fiber.Ctx) error {
		username := c.Params("username")

		uname, err := url.PathUnescape(username)
		if err != nil || uname == "" {
			return c.
				Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeInvalidUsername,
					Message: uerrors.MessageInvalidUsername,
				})
		}

		user, err := usersDB.GetUserByUsername(username)
		if err != nil {
			code := err.(*uerrors.Error).Code

			return c.
				Status(uerrors.ToHTTPStatusCode(code)).
				JSON(err)
		}

		return c.JSON(user)
	})

	users.Get("/id/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		id, err := url.PathUnescape(id)
		if err != nil || id == "" {
			return c.
				Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeInvalidUserID,
					Message: uerrors.MessageInvalidUserID,
				})
		}

		uid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return c.
				Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeInvalidUserID,
					Message: uerrors.MessageInvalidUserID,
				})
		}

		user, err := usersDB.GetUserByID(uid)
		if err != nil {
			code := err.(*uerrors.Error).Code

			return c.
				Status(uerrors.ToHTTPStatusCode(code)).
				JSON(err)
		}

		return c.JSON(user)
	})

	users.Post("/", func(c *fiber.Ctx) error {
		c.Accepts(fiber.MIMEApplicationJSON)

		var newUser api.User

		if len(c.Body()) == 0 {
			return c.
				Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeEmptyBody,
					Message: uerrors.MessageEmptyBody,
				})
		}

		if err := json.Unmarshal(c.Body(), &newUser); err != nil {
			return c.
				Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeInvalidUserPost,
					Message: fmt.Sprintf("%s %s", uerrors.MessageInvalidUserPost, err.Error()),
				})
		}

		createdUser, err := usersDB.CreateUser(&newUser)
		if err != nil {
			return c.
				Status(uerrors.ToHTTPStatusCode(err.(*uerrors.Error).Code)).
				JSON(err)
		}

		return c.
			Status(fiber.StatusCreated).
			JSON(createdUser)
	})

	users.Put("/:id", func(c *fiber.Ctx) error {
		c.Accepts(fiber.MIMEApplicationJSON)

		var userToUpd api.User

		if len(c.Body()) == 0 {
			return c.
				Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeEmptyBody,
					Message: uerrors.MessageEmptyBody,
				})
		}

		if err := json.Unmarshal(c.Body(), &userToUpd); err != nil {
			return c.
				Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeInvalidUserPost,
					Message: fmt.Sprintf("%s %s", uerrors.MessageInvalidUserPost, err.Error()),
				})
		}

		id := c.Params("id")

		id, err := url.PathUnescape(id)
		if err != nil || id == "" {
			return c.
				Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeInvalidUserID,
					Message: uerrors.MessageInvalidUserID,
				})
		}

		uid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return c.
				Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeInvalidUserID,
					Message: uerrors.MessageInvalidUserID,
				})
		}

		if userToUpd.ID != uid {
			// Just in case...
			userToUpd.ID = uid
		}

		// TODO: check if user is admin or owner of this profile
		if err = usersDB.UpdateUser(uid, &userToUpd); err != nil {
			return c.
				Status(uerrors.ToHTTPStatusCode(err.(*uerrors.Error).Code)).
				JSON(err)
		}

		return c.SendStatus(fiber.StatusOK)
	})

	users.Delete("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		id, err := url.PathUnescape(id)
		if err != nil || id == "" {
			return c.
				Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeInvalidUserID,
					Message: uerrors.MessageInvalidUserID,
				})
		}

		uid, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return c.
				Status(fiber.StatusBadRequest).
				JSON(&uerrors.Error{
					Code:    uerrors.CodeInvalidUserID,
					Message: uerrors.MessageInvalidUserID,
				})
		}

		hardDelete, err := url.QueryUnescape(strings.ToLower(c.Query("hard_delete", "false")))
		if err != nil {
			hardDelete = "false"
		}

		// TODO: check if user CAN hard delete
		// TODO: check if user is admin or owner of this profile

		if err = usersDB.DeleteUser(uid, hardDelete == "true"); err != nil {
			return c.
				Status(uerrors.ToHTTPStatusCode(err.(*uerrors.Error).Code)).
				JSON(err)
		}

		return c.SendStatus(fiber.StatusGone)
	})

	internalEndpoints := fiber.New(fiber.Config{
		AppName:               fiberAppName,
		ReadTimeout:           time.Minute,
		DisableStartupMessage: verbosity > 0,
	})

	internalEndpoints.Get("/readyz", func(c *fiber.Ctx) error {
		if db != nil {
			return c.SendStatus(fiber.StatusOK)
		}

		return c.SendStatus(fiber.StatusInternalServerError)
	})

	internalEndpoints.Get("/livez", func(c *fiber.Ctx) error {
		return c.SendStatus(fiber.StatusOK)
	})

	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Err(err).Msg("error while listening")
		}
	}()

	go func() {
		if err := internalEndpoints.Listen(":8081"); err != nil {
			log.Err(err).Msg("error while listening")
		}
	}()

	// Graceful Shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Info().Msg("shutting down...")
	if err := app.Shutdown(); err != nil {
		log.Err(err).Msg("error while waiting for server to shutdown")
	}
	if err := internalEndpoints.Shutdown(); err != nil {
		log.Err(err).Msg("error while waiting for server to shutdown")
	}
	log.Info().Msg("goodbye!")
}
