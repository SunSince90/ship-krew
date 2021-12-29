package main

import (
	"flag"
	"net/url"
	"os"
	"os/signal"
	"time"

	udb "github.com/SunSince90/ship-krew/users/api/internal/database"
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

	go func() {
		if err := app.Listen(":8080"); err != nil {
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
	log.Info().Msg("goodbye!")
}
