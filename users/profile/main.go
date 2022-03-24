package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path"
	"time"

	"github.com/asimpleidea/ship-krew/users/api/pkg/api"
	uerrors "github.com/asimpleidea/ship-krew/users/api/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/rs/zerolog"
)

const (
	fiberAppName          string        = "Profile Backend"
	defaultApiTimeout     time.Duration = time.Minute
	defaultViewsDirectory string        = "/views"
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
		verbosity      int
		usersApiAddr   string
		timeout        time.Duration
		viewsDirectory string
		appViews       string
	)

	flag.IntVar(&verbosity, "verbosity", 1, "the verbosity level")

	// TODO: https, not http
	flag.StringVar(&usersApiAddr, "users-api-address", "http://users-api", "the address of the users server API")
	flag.DurationVar(&timeout, "timeout", 2*time.Minute, "requests timeout")
	flag.StringVar(&viewsDirectory, "views-directory", defaultViewsDirectory,
		"Directory containing views.")
	flag.Parse()

	log = zerolog.New(os.Stderr).With().Logger()
	log.Info().Int("verbosity", verbosity).Msg("starting...")

	{
		logLevels := [4]zerolog.Level{zerolog.DebugLevel, zerolog.InfoLevel, zerolog.ErrorLevel}
		log = log.Level(logLevels[verbosity])
	}

	viewsDir := path.Join(viewsDirectory, "public")
	appViews = path.Join("apps", "profile")

	// TODO: if not available should fail
	engine := html.New(viewsDir, ".html")

	// TODO: authenticate to users server with APIKey

	app := fiber.New(fiber.Config{
		AppName:               fiberAppName,
		ReadTimeout:           time.Minute,
		DisableStartupMessage: verbosity > 0,
		Views:                 engine,
	})

	app.Get("/profiles/:username", func(c *fiber.Ctx) error {
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

		return c.Render(path.Join(appViews, "index"), fiber.Map{
			"Title": fmt.Sprintf("Hello, %s!", user.DisplayName),
			// TODO: find a way to do this in a better way, maybe from template?
			"EditURL": path.Join("u", user.Username, "edit"),
			"User":    user,
		})
	})

	app.Get("/profiles/:username/edit", func(c *fiber.Ctx) error {
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

		return c.Render(path.Join(appViews, "edit_profile"), fiber.Map{
			"User": user,
		})
	})

	app.Post("/profiles/:username/edit", func(c *fiber.Ctx) error {
		// TODO:
		// - Check if you can do this (OPA) or maybe let the API do this?
		// - Handle case in which this is called via AJAX

		ctx, canc := context.WithTimeout(context.Background(), defaultApiTimeout)
		usr, err := getUserByUsername(ctx, usersApiAddr, c.Params("username"))
		if err != nil {
			// TODO: parse the error and return an html of the error, not
			// simple text.
			canc()
			return c.SendStatus(fiber.StatusInternalServerError)
		}
		canc()

		usrToUpdate := usr.Clone()

		const (
			formUsername    = "edit_username"
			formDisplayName = "edit_display_name"
			formBio         = "edit_bio"
		)

		{
			// TODO:
			// - validation and check for prohibited words
			// - check from settings how many times you can change it in x days.
			editedUsername := c.FormValue(formUsername)
			if editedUsername != "" && editedUsername != usr.Username {
				usrToUpdate.Username = editedUsername
			}
		}

		{
			// TODO:
			// - validation and check for prohibited words
			editedDisplayName := c.FormValue(formDisplayName)
			if editedDisplayName != "" && editedDisplayName != usr.DisplayName {
				usrToUpdate.DisplayName = editedDisplayName
			}
		}

		{
			// TODO:
			// - validation
			editedBio := c.FormValue(formBio)
			if editedBio != "" {
				if usrToUpdate.Bio == nil || (usrToUpdate.Bio != nil && editedBio != *usrToUpdate.Bio) {
					usrToUpdate.Bio = &editedBio
				}
			}
		}

		ctx, canc = context.WithTimeout(context.Background(), defaultApiTimeout)
		defer canc()
		if err := updateUser(ctx, usersApiAddr, usrToUpdate); err != nil {
			// TODO:
			// - Parse the error and decide what to do
			// - Send json if ajax or html if not
			var e *uerrors.Error
			if errors.As(err, &e) {
				return c.Status(uerrors.ToHTTPStatusCode(e.Code)).
					JSON(e)
			}

			return c.Status(fiber.StatusInternalServerError).
				Send([]byte(err.Error()))
		}

		return c.SendStatus(fiber.StatusOK)
	})

	// TODO: only do readiness probe.
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

// TODO: this must be integrated in the client
func updateUser(ctx context.Context, usersApiAddr string, user *api.User) error {
	usrBody, err := json.Marshal(user)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPut,
		fmt.Sprintf("%s/users/%d", usersApiAddr, user.ID),
		bytes.NewReader(usrBody))
	if err != nil {
		return err
	}

	// TODO: use cookies in client?
	cl := &http.Client{}
	resp, err := cl.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// TODO: handle errors better
	if resp.StatusCode != fiber.StatusOK {
		var oerr uerrors.Error
		if err := json.NewDecoder(resp.Body).Decode(&err); err != nil {
			return &oerr
		}

		return &oerr
	}

	return nil
}
