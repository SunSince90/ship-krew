package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/rs/zerolog"
)

const (
	fiberAppName string = "Profile Backend"
)

var (
	log zerolog.Logger
)

type Elems struct {
	Color string
	Val   string
}

func main() {
	var (
		verbosity    int
		usersApiAddr string
		timeout      time.Duration
	)

	flag.IntVar(&verbosity, "verbosity", 1, "the verbosity level")

	// TODO: https, not http
	flag.StringVar(&usersApiAddr, "users-api-address", "http://users-api", "the address of the users server API")
	flag.DurationVar(&timeout, "timeout", 2*time.Minute, "requests timeout")
	flag.Parse()

	log = zerolog.New(os.Stderr).With().Logger()
	log.Info().Int("verbosity", verbosity).Msg("starting...")

	{
		logLevels := [4]zerolog.Level{zerolog.DebugLevel, zerolog.InfoLevel, zerolog.ErrorLevel}
		log = log.Level(logLevels[verbosity])
	}

	engine := html.New("./views", ".html")

	// TODO: authenticate to users server with APIKey

	app := fiber.New(fiber.Config{
		AppName:               fiberAppName,
		ReadTimeout:           time.Minute,
		DisableStartupMessage: verbosity > 0,
		Views:                 engine,
	})

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title":  "Hello, World!",
			"Things": []Elems{{Color: "red", Val: "one"}, {Color: "blue", Val: "two"}},
		})
	})

	app.Get("/u/:username", func(c *fiber.Ctx) error {
		// TODO: get API users username

		return c.Render("index", fiber.Map{
			"Title":  fmt.Sprintf("Hello, %s!", c.Params("username")),
			"Things": []Elems{{Color: "red", Val: "one"}, {Color: "blue", Val: "two"}},
		})
	})

	app.Get("/u/:username/edit", func(c *fiber.Ctx) error {
		// TODO: get API users username

		return c.Render("edit_profile", fiber.Map{
			"DisplayName": "TODO",
			"Username":    c.Params("username"),
		})
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
