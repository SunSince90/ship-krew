package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/asimpleidea/ship-krew/users/api/pkg/api"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/rs/zerolog"
)

const (
	fiberAppName      string        = "Profile Backend"
	defaultApiTimeout time.Duration = time.Minute
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
		// TODO: should username be sanitized?
		ctx, canc := context.WithTimeout(context.Background(), defaultApiTimeout)
		user, err := getUserByUsername(ctx, usersApiAddr, c.Params("username"))
		if err != nil {
			// TODO: parse the erorr and return an html of the error, not
			// simple text.
			canc()
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		canc()

		return c.Render("index", fiber.Map{
			"Title": fmt.Sprintf("Hello, %s!", user.DisplayName),
			// TODO: find a way to do this in a better way, maybe from template?
			"EditURL": path.Join("u", user.Username, "edit"),
			"User":    user,
		})
	})

	app.Get("/u/:username/edit", func(c *fiber.Ctx) error {
		// TODO: get API users username
		ctx, canc := context.WithTimeout(context.Background(), defaultApiTimeout)
		user, err := getUserByUsername(ctx, usersApiAddr, c.Params("username"))
		if err != nil {
			// TODO: parse the erorr and return an html of the error, not
			// simple text.
			canc()
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		canc()

		return c.Render("edit_profile", fiber.Map{
			"User": user,
		})
	})

	app.Post("/u/:username/edit", func(c *fiber.Ctx) error {
		// TODO: send to UPDATE api
		// TODO: check if user can actually do this or leave this to api?
		return c.SendString("editing")
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

// TODO: this must be integrated in the client
func getUserByUsername(ctx context.Context, usersApiAddr, username string) (*api.User, error) {
	req, err := http.NewRequestWithContext(ctx,
		http.MethodGet,
		fmt.Sprintf("%s/users/username/%s", usersApiAddr, username),
		nil)
	if err != nil {
		return nil, err
	}

	// TODO: use cookies in client?
	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var user api.User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}

	return &user, nil
}
