package main

import (
	"os"
	"time"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	flag "github.com/spf13/pflag"
)

var (
	log zerolog.Logger
)

type options struct {
	verbosity int
}

func main() {
	verbosityLevels := []zerolog.Level{zerolog.DebugLevel, zerolog.InfoLevel, zerolog.ErrorLevel, zerolog.FatalLevel}
	log = zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = time.RFC3339

	opts := &options{}
	flag.IntVarP(&opts.verbosity, "verbosity", "v", 1, "the log verbosity level: 0 is the most verbose and 3 the quietest.")
	flag.Parse()

	if opts.verbosity < 0 || opts.verbosity > 3 {
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		log.Info().
			Int("provided-verbosity", opts.verbosity).
			Int("default-level", 1).
			Msg("invalid verbosity level provided: using default level")
	} else {
		zerolog.SetGlobalLevel(verbosityLevels[opts.verbosity])
	}

	log.Info().Msg("starting...")

	// Set up the routes
	app := fiber.New()
	api := app.Group("/api", func(c *fiber.Ctx) error {
		return c.Next()
	})
	api.Get("/users", func(c *fiber.Ctx) error {
		log.Info().Msg("called /users")
		return c.SendString("called /users")
	})

	canceling := make(chan struct{})
	go func() {
		app.Listen(":8080")
		close(canceling)
	}()

	<-canceling
}
