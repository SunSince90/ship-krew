package main

import (
	"flag"
	"os"
	"time"

	"github.com/SunSince90/ship-krew/users/api/pkg/database"
	"github.com/rs/zerolog"
)

var (
	log zerolog.Logger
)

func main() {
	var (
		verbosity  int
		dbSettings = &database.Settings{}
	)

	flag.IntVar(&verbosity, "verbosity", 0, "the verbosity level")

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

	// TODO: continue
}
