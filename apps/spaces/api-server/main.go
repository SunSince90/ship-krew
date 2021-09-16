package main

import (
	"os"

	"github.com/rs/zerolog"
	flag "github.com/spf13/pflag"
)

var (
	log zerolog.Logger
)

func main() {
	usersApiAddress := "users-api-server.ship-krew-backend"
	testSpaces := 0
	verbosity := 1

	flag.StringVar(&usersApiAddress, "users-api-server", usersApiAddress, "the address where to contact the users API server")
	flag.IntVar(&testSpaces, "test-spaces", 0, "the number of test spaces to create")
	flag.IntVarP(&verbosity, "verbosity", "v", 1, "the verbosity level")
	flag.Parse()

	log = zerolog.New(os.Stderr).With().Timestamp().Logger()

	verbosityLevels := []zerolog.Level{zerolog.DebugLevel, zerolog.InfoLevel, zerolog.ErrorLevel}
	if verbosity < 0 || verbosity > len(verbosityLevels) {
		log.Error().
			Int("verbosity", verbosity).
			Int("default", 1).
			Msg("invalid verbosity level provided, reverting to default...")
		verbosity = 1
	}

	if testSpaces > 0 {
		log.Info().
			Int("test-spaces", testSpaces).
			Msg("test spaces requested")
		verbosity = 0
	}

	log = log.Level(verbosityLevels[verbosity]).With().Logger()
	log.Info().Msg("starting...")
}
