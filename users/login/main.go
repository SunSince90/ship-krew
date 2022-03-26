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
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path"
	"strings"
	"time"

	"github.com/asimpleidea/ship-krew/users/api/pkg/api"
	uerrors "github.com/asimpleidea/ship-krew/users/api/pkg/errors"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/encryptcookie"
	"github.com/gofiber/template/html"
	"github.com/rs/zerolog"
	"gopkg.in/yaml.v3"
)

// TODO: Follow https://stackoverflow.com/questions/244882/what-is-the-best-way-to-implement-remember-me-for-a-website
// TODO: follow https://paragonie.com/blog/2015/04/secure-authentication-php-with-long-term-persistence#title.2

const (
	fiberAppName          string        = "Login Backend"
	defaultApiTimeout     time.Duration = time.Minute
	defaultPongTimeout    time.Duration = 30 * time.Second
	defaultViewsDirectory string        = "/views"
)

var (
	log zerolog.Logger
)

func main() {
	var (
		verbosity      int
		usersApiAddr   string
		timeout        time.Duration
		cookieKey      string
		redisEndpoint  string
		redisPassword  string
		viewsDirectory string
		appViews       string
	)

	flag.IntVar(&verbosity, "verbosity", 1, "the verbosity level")

	// TODO: https, not http
	flag.StringVar(&usersApiAddr, "users-api-address", "http://users-api", "the address of the users server API")
	flag.DurationVar(&timeout, "timeout", 2*time.Minute, "requests timeout")

	// TODO: this should be pulled from secrets
	flag.StringVar(&cookieKey, "cookie-key", "", "The key to un-encrypt cookies")

	// TODO:
	// - use default
	// - use secrets for authentication
	// - what is a good default for this?
	flag.StringVar(&redisEndpoint, "redis-endpoints", "http://localhost:6379",
		"Endpoints where to contact redis.")
	// TODO: this must be a certificate when stable.
	flag.StringVar(&redisPassword, "redis-password", "",
		"Authentication password for redis.")

	flag.StringVar(&viewsDirectory, "views-directory", defaultViewsDirectory,
		"Root directory containing views.")
	flag.Parse()

	log = zerolog.New(os.Stderr).With().Logger()
	log.Info().Int("verbosity", verbosity).Msg("starting...")

	{
		logLevels := [4]zerolog.Level{zerolog.DebugLevel, zerolog.InfoLevel, zerolog.ErrorLevel}
		log = log.Level(logLevels[verbosity])
	}

	if cookieKey == "" {
		log.Fatal().Err(errors.New("no cookie key set")).Msg("fatal error occurred")
	}

	// ------------------------------------
	// Get the client from redis (for sessions)
	// ------------------------------------

	sessClient := redis.NewClient(&redis.Options{
		Addr:     redisEndpoint,
		Password: redisPassword,
		// TODO: define the database from flags.
		DB: 0,
	})
	defer sessClient.Close()

	if err := func() error {
		ctx, canc := context.WithTimeout(context.TODO(), defaultPongTimeout)
		defer canc()

		_, err := sessClient.Ping(ctx).Result()
		return err
	}(); err != nil {
		log.Fatal().Err(err).Msg("could not connect to redis")
		return // unnecessary but for readability
	}

	viewsDir := path.Join(viewsDirectory, "public")
	appViews = path.Join("apps", "login")

	// TODO: if not available should fail
	engine := html.New(viewsDir, ".html")

	// TODO: authenticate to users server with APIKey

	app := fiber.New(fiber.Config{
		AppName:               fiberAppName,
		ReadTimeout:           time.Minute,
		DisableStartupMessage: verbosity > 0,
		Views:                 engine,
	})

	app.Use(encryptcookie.New(encryptcookie.Config{
		Key: cookieKey,
	}))

	app.Get("/login", func(c *fiber.Ctx) error {
		// TODO: make this whole function better
		ctx, canc := context.WithTimeout(context.Background(), defaultApiTimeout)
		usrSession, sessionID, _ := getSessionFromRedis(ctx, c, sessClient)
		// TODO: check if error is from redis, if so internal server error
		// otherwise just login
		canc()
		if usrSession != nil {
			if !usrSession.Expired() {
				if usrSession.DaysTillExpiration() < 3 {
					crCtx, crCanc := context.WithTimeout(context.Background(), defaultApiTimeout)
					usrSession.Expiration = time.Now().AddDate(0, 0, 7)
					err := createSessionOnRedis(crCtx, sessClient, *sessionID, usrSession)
					crCanc()
					if err != nil {
						log.Err(err).Str("session-id", *sessionID).
							Int64("user-id", usrSession.UserID).
							Msg("error while trying to update session")
					}

				}

				return c.Status(fiber.StatusNotFound).SendString("already logged in")
			}

			func() {
				delCtx, delCanc := context.WithTimeout(context.Background(), defaultApiTimeout)
				defer delCanc()
				if err := deleteSession(delCtx, c, sessClient, *sessionID); err != nil {
					log.Err(err).Str("session-id", *sessionID).
						Int64("user-id", usrSession.UserID).
						Msg("error while trying to delete session")
				}
			}()
		}

		// TODO:
		// - check if session exists on Redis
		// - check if it corresponds to this user
		// - check if not expired

		// TODO:
		// - This must be called login
		return c.Render(path.Join(appViews, "index"), fiber.Map{
			"Title": "Login",
		})
	})

	app.Post("/login", func(c *fiber.Ctx) error {
		ctx, canc := context.WithTimeout(context.Background(), defaultApiTimeout)
		usrSession, sessionID, _ := getSessionFromRedis(ctx, c, sessClient)
		// TODO: check if error is from Redis, if so internal server error
		// otherwise just login
		canc()
		if usrSession != nil {
			if !usrSession.Expired() {
				if usrSession.DaysTillExpiration() < 3 {
					crCtx, crCanc := context.WithTimeout(context.Background(), defaultApiTimeout)
					usrSession.Expiration = time.Now().AddDate(0, 0, 7)
					err := createSessionOnRedis(crCtx, sessClient, *sessionID, usrSession)
					crCanc()
					if err != nil {
						log.Err(err).Str("session-id", *sessionID).
							Int64("user-id", usrSession.UserID).
							Msg("error while trying to update session")
					}
				}

				return c.Status(fiber.StatusNotFound).SendString("already logged in")
			}

			func() {
				delCtx, delCanc := context.WithTimeout(context.Background(), defaultApiTimeout)
				defer delCanc()
				if err := deleteSession(delCtx, c, sessClient, *sessionID); err != nil {
					log.Err(err).Str("session-id", *sessionID).
						Int64("user-id", usrSession.UserID).
						Msg("error while trying to delete session")
				}
			}()
		}

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

		ctx, canc = context.WithTimeout(context.Background(), defaultApiTimeout)
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

			return c.Status(fiber.StatusBadRequest).SendString("not ok")
		}
		canc()

		if passwordIsCorrect(pwd, usr.Base64PasswordHash, usr.Base64Salt) {
			fmt.Println("password is correct")
			// TODO: generate a good session ID
			sessionID := "testing"
			c.Cookie(&fiber.Cookie{
				Name:  "session",
				Value: sessionID,
			})

			ctx, canc = context.WithTimeout(context.Background(), defaultApiTimeout)
			defer canc()
			if err := createSessionOnRedis(ctx, sessClient, sessionID, &UserSession{
				CreatedAt:  time.Now(),
				UserID:     usr.ID,
				Expiration: time.Now().AddDate(0, 0, 7),
			}); err != nil {
				return c.Status(fiber.StatusInternalServerError).
					Send([]byte(err.Error()))
			}

			return c.Status(fiber.StatusOK).Send([]byte("ok"))
		}

		// TODO: cookie

		return c.Status(fiber.StatusOK).
			Send([]byte("does not match"))
	})

	app.Get("/logout", func(c *fiber.Ctx) error {
		ctx, canc := context.WithTimeout(context.Background(), defaultApiTimeout)
		usrSession, sessID, _ := getSessionFromRedis(ctx, c, sessClient)
		canc()

		if usrSession == nil {
			return c.Status(fiber.StatusNotFound).
				SendString("not logged int")
		}

		// TODO:
		// - check if session exists on Redis
		// - check if it corresponds to this user
		// - check if not expired

		c.ClearCookie("session")
		ctx, canc = context.WithTimeout(context.Background(), defaultApiTimeout)
		if err := deleteSession(ctx, c, sessClient, *sessID); err != nil {
			log.Err(err).Str("session-id", *sessID).
				Int64("user-id", usrSession.UserID).
				Msg("error while trying to delete session")
			canc()
		}
		canc()

		// TODO: redirect
		return c.Status(fiber.StatusOK).SendString("ok")
	})

	app.Get("/signup", func(c *fiber.Ctx) error {
		ctx, canc := context.WithTimeout(context.Background(), defaultApiTimeout)
		usrSession, sessionID, _ := getSessionFromRedis(ctx, c, sessClient)
		canc()

		if usrSession != nil {
			if usrSession.Expired() {
				func() {
					delCtx, delCanc := context.WithTimeout(context.Background(), defaultApiTimeout)
					defer delCanc()
					if err := deleteSession(delCtx, c, sessClient, *sessionID); err != nil {
						log.Err(err).Str("session-id", *sessionID).
							Int64("user-id", usrSession.UserID).
							Msg("error while trying to delete session")
					}
				}()
			} else {
				return c.Status(fiber.StatusNotFound).SendString("already logged in")
			}
		}

		return c.Render(path.Join(appViews, "signup"), fiber.Map{
			"Title": "Signup",
		})
	})

	app.Post("/signup", func(c *fiber.Ctx) error {
		// TODO:
		// - validate form values
		// - check if email already exists

		{
			ctx, canc := context.WithTimeout(context.Background(), defaultApiTimeout)
			usr, err := getUserByUsername(ctx, usersApiAddr, c.FormValue("signup_username"))
			canc()
			if err != nil {
				var e *uerrors.Error
				if errors.As(err, &e) && e.Code != uerrors.CodeUserNotFound {
					// TODO:
					// - html
					// - better parsing
					return c.Status(fiber.StatusInternalServerError).SendString(e.Error())
				}
			}
			if usr != nil {
				// TODO:
				// is good code?
				// html
				return c.Status(fiber.StatusBadRequest).
					SendString("a user with this username already exists")
			}
		}

		{
			if c.FormValue("signup_password") != c.FormValue("signup_confirm_password") {
				return c.Status(fiber.StatusBadGateway).SendString("passwords do not match")
			}
		}

		userToCreate := &api.User{
			Username:    c.FormValue("signup_username"),
			DisplayName: c.FormValue("signup_username"),
			Email: func() *string {
				email := c.FormValue("signup_email")
				return &email
			}(),
			RegistrationIP: func() *net.IP {
				// TODO: make this better
				endpoint := c.Context().RemoteAddr().String()
				colon := strings.Index(endpoint, ":")

				ip := net.ParseIP(endpoint[0:colon])
				return &ip
			}(),
		}

		{
			passBytes := bytes.NewBufferString(c.FormValue("signup_password")).Bytes()
			// TODO: this should not use sha256!
			passHash := sha256.Sum256(passBytes)
			base64Pass := base64.StdEncoding.EncodeToString(passHash[:])
			userToCreate.Base64PasswordHash = &base64Pass
		}

		{
			ctx, canc := context.WithTimeout(context.Background(), defaultApiTimeout)
			if err := createUser(ctx, usersApiAddr, userToCreate); err != nil {
				canc()
				var e *uerrors.Error
				if errors.As(err, &e) {
					// TODO:
					// - html
					// - better parsing
					return c.Status(uerrors.ToHTTPStatusCode(e.Code)).
						JSON(e)
				}

				return c.Status(fiber.StatusInternalServerError).SendString(e.Error())
			}
			canc()
		}

		return c.Status(fiber.StatusOK).SendString("ok")
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

	// TODO: better way to handle these internal server error

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeInternalServerError,
			Message: uerrors.MessageInternalServerError,
		}
	}

	if resp.StatusCode != fiber.StatusOK {
		var e uerrors.Error
		if err := json.Unmarshal(body, &e); err != nil {
			return nil, &uerrors.Error{
				Code:    uerrors.CodeInternalServerError,
				Message: uerrors.MessageInternalServerError,
			}
		}

		return nil, &e
	}

	var user api.User
	if err := json.Unmarshal(body, &user); err != nil {
		return nil, &uerrors.Error{
			Code:    uerrors.CodeInternalServerError,
			Message: uerrors.MessageInternalServerError,
		}
	}

	return &user, nil
}

func createUser(ctx context.Context, usersApiAddr string, usr *api.User) error {
	bodyToSend, err := json.Marshal(usr)
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx,
		http.MethodPost,
		fmt.Sprintf("%s/users", usersApiAddr),
		bytes.NewBuffer(bodyToSend))
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

	// TODO: better way to handle these internal server error

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return &uerrors.Error{
			Code:    uerrors.CodeInternalServerError,
			Message: uerrors.MessageInternalServerError,
		}
	}

	if resp.StatusCode >= 300 {
		var e uerrors.Error
		if err := json.Unmarshal(body, &e); err != nil {
			return &uerrors.Error{
				Code:    uerrors.CodeInternalServerError,
				Message: uerrors.MessageInternalServerError,
			}
		}

		return &e
	}

	return nil
}

// TODO: this may need to be better and maybe done on client
func passwordIsCorrect(provided string, expected, salt *string) bool {
	decodedExpected, _ := base64.URLEncoding.DecodeString(*expected)
	decodedSalt, _ := base64.URLEncoding.DecodeString(*salt)

	digestProvided := sha256.Sum256([]byte(provided))
	passWithSalt := append(digestProvided[:], decodedSalt...)

	return bytes.Equal(passWithSalt, decodedExpected)
}

type UserSession struct {
	CreatedAt  time.Time `json:"created_at" yaml:"createdAt"`
	UserID     int64     `json:"user_id" yaml:"userId"`
	Expiration time.Time `json:"expiration" yaml:"expiration"`
}

func (u *UserSession) Expired() bool {
	return time.Now().After(u.Expiration)
}

func (u *UserSession) DaysTillExpiration() int {
	if u.Expired() {
		return -1
	}

	return int(time.Until(u.Expiration).Hours() / 24)
}

func getSessionFromRedis(ctx context.Context, fctx *fiber.Ctx, sessionsClient *redis.Client) (*UserSession, *string, error) {
	// TODO: not sure about bringing fctx here
	sessionID := fctx.Cookies("session", "")
	if sessionID == "" {
		return nil, nil, nil
	}

	val, err := sessionsClient.Get(ctx, path.Join("sessions", sessionID)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, nil, nil
		}

		return nil, nil, fmt.Errorf("error while getting session key: %w", err)
	}

	var usrSession UserSession
	if err := yaml.NewDecoder(strings.NewReader(val)).Decode(&usrSession); err != nil {
		return nil, nil, fmt.Errorf("error while unmarshalling value from value: %w", err)
	}

	return &usrSession, &sessionID, nil
}

func createSessionOnRedis(ctx context.Context, sessionClient *redis.Client, sessionID string, usrSession *UserSession) error {
	val, err := yaml.Marshal(usrSession)
	if err != nil {
		return fmt.Errorf("could not marshal data: %w", err)
	}

	return sessionClient.
		Set(ctx, path.Join("sessions", sessionID), val, time.Until(usrSession.Expiration)).
		Err()
}

func deleteSession(ctx context.Context, fctx *fiber.Ctx, sessionClient *redis.Client, sessionID string) error {
	fctx.ClearCookie("session")

	return sessionClient.Del(ctx, path.Join("sessions", sessionID)).Err()
}
