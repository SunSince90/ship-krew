package main

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/SunSince90/ship-krew/users/api-server/pkg/api"
	fakeit "github.com/brianvoe/gofakeit/v6"
	fiber "github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	uuid "github.com/satori/go.uuid"
	flag "github.com/spf13/pflag"
)

var (
	log zerolog.Logger
)

type options struct {
	verbosity int
}

func main() {
	testUsers := 0
	verbosityLevels := []zerolog.Level{zerolog.DebugLevel, zerolog.InfoLevel, zerolog.ErrorLevel, zerolog.FatalLevel}
	log = zerolog.New(os.Stderr).With().Timestamp().Logger()
	zerolog.TimeFieldFormat = time.RFC3339

	opts := &options{}
	flag.IntVarP(&opts.verbosity, "verbosity", "v", 1, "the log verbosity level: 0 is the most verbose and 3 the quietest.")
	flag.IntVar(&testUsers, "test-users", 0, "the number of test users to create. If more than 1, then test mode will enabled and verbosity will be put to 0 automatically")
	flag.Parse()

	if testUsers > 0 {
		opts.verbosity = 0
		log.Info().Int("test-users", testUsers).Msg("test mode requested")
	}

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

	if testUsers > 0 {
		createTestUsers(testUsers)
	}

	// Set up the routes
	app := fiber.New()
	api := app.Group("/api", func(c *fiber.Ctx) error {
		return c.Next()
	})
	api.Get("/users", listUsersHandler)
	api.Get("/users/version", getVersionHandler)
	api.Post("/users", createUserHandler)
	api.Get("/users/:name", getUserHandler)
	api.Put("/users/:name", updateUserHandler)
	api.Delete("/users/:name", deleteUserHandler)

	// Probes
	probes := fiber.New()
	probes.Get("/healthz", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})
	probes.Get("/ready", func(c *fiber.Ctx) error {
		return c.SendStatus(200)
	})

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		app.Listen(":8080")
	}()
	go func() {
		defer wg.Done()
		probes.Listen(":8081")
	}()

	wg.Wait()
}

// TODO: this is just for testing and will be removed.
var usersList map[string]api.User

func createTestUsers(count int) {
	log.Info().Msg("creating test users...")
	fmt.Println()

	if len(usersList) == 0 {
		usersList = map[string]api.User{}
	}

	for i := 0; i < count; i++ {
		userToCreate := api.User{
			ID:        uuid.NewV4().String(),
			Name:      fakeit.Username(),
			Bio:       fakeit.Quote(),
			CreatedAt: fakeit.Date(),
		}

		fmt.Println(userToCreate)
		usersList[userToCreate.Name] = userToCreate
	}
	fmt.Println()
}

func updateTestUser(name string, newData api.User) {
	usersList[name] = newData
}

func deleteTestUser(name string) {
	delete(usersList, name)
}

var (
	// version   string
	gitCommit string
)

func getVersionHandler(c *fiber.Ctx) error {
	return c.JSON(map[string]string{
		// "Version":    version,
		"Git Commit": gitCommit,
	})
}
