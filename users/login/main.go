package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/asimpleidea/ship-krew/users/api/pkg/api"
	uerrors "github.com/asimpleidea/ship-krew/users/api/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"
	"github.com/rs/zerolog"
)

const (
	fiberAppName      string        = "Login Backend"
	defaultApiTimeout time.Duration = time.Minute
)

var (
	log zerolog.Logger
)

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
		// TODO:
		// - This must be called login
		return c.Render("index", fiber.Map{
			"Title": "Login",
		})
	})

	app.Post("/", func(c *fiber.Ctx) error {
		// TODO:
		// - validations
		// - check if values are actually provided
		// - check if ajax
		const (
			formUsername = "login_username"
			formPassword = "login_password"
		)

		username := c.FormValue(formUsername)
		pwd := c.FormValue(formPassword)

		ctx, canc := context.WithTimeout(context.Background(), defaultApiTimeout)
		usr, err := getUserByUsername(ctx, usersApiAddr, username)
		if err != nil {
			canc()
			// TODO:
			// - do not disclose that this user does not exist, but just say
			// 	that this user-pwd combination was not found
			var e *uerrors.Error
			if errors.As(err, &e) {
				return c.Status(uerrors.ToHTTPStatusCode(e.Code)).
					JSON(e)
			}
		}
		canc()

		if passwordIsCorrect(pwd, usr.Base64PasswordHash, usr.Base64Salt) {
			return c.Status(fiber.StatusOK).Send([]byte("ok"))
		}

		// TODO: cookie

		return c.Status(fiber.StatusOK).
			Send([]byte("does not match"))
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

// TODO: this may need to be better and maybe done on client
func passwordIsCorrect(provided string, expected, salt *string) bool {
	decodedExpected, _ := base64.URLEncoding.DecodeString(*expected)
	decodedSalt, _ := base64.URLEncoding.DecodeString(*salt)

	digestProvided := sha256.Sum256([]byte(provided))
	passWithSalt := append(digestProvided[:], decodedSalt...)

	return bytes.Compare(passWithSalt, decodedExpected) == 0
}
