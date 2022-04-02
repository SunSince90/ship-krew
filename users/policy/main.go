package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/rs/zerolog"
)

const (
	defaultRegoDirectory string        = "/rego"
	defaultApiTimeout    time.Duration = time.Minute
	defaultPongTimeout   time.Duration = 30 * time.Second
)

var (
	log zerolog.Logger
)

func main() {
	var (
		verbosity     int
		regoDirectory string
	)

	flag.IntVar(&verbosity, "verbosity", 1, "the verbosity level")

	flag.StringVar(&regoDirectory, "rego-directory", defaultRegoDirectory,
		"Root directory containing rego files.")
	flag.Parse()

	log = zerolog.New(os.Stderr).With().Logger()
	log.Info().Int("verbosity", verbosity).Msg("starting...")

	{
		logLevels := [4]zerolog.Level{zerolog.DebugLevel, zerolog.InfoLevel, zerolog.ErrorLevel}
		log = log.Level(logLevels[verbosity])
	}

	// --------------------------------------------
	// Load rego files from directory
	// --------------------------------------------

	// TODO...

	// --------------------------------------------
	// Start the gRPC server
	// --------------------------------------------

	// TODO...

	// Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop

	log.Info().Msg("shutting down...")
	log.Info().Msg("goodbye!")
}
