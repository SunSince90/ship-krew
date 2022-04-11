package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/open-policy-agent/opa/ast"
	"github.com/open-policy-agent/opa/rego"
	"github.com/open-policy-agent/opa/util"
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

	ctx, canc := context.WithCancel(context.Background())
	_ = canc

	// --------------------------------------------
	// Load rego files from directory
	// --------------------------------------------

	ver, err := newVerifier(ctx, regoDirectory)
	if err != nil {
		log.Fatal().Err(err).Msg("could not load verifier")
		return
	}

	// TODO: use the verifier
	_ = ver

	// --------------------------------------------
	// Start the gRPC server
	// --------------------------------------------

	// TODO...

	// Graceful Shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	<-stop
	canc()

	log.Info().Msg("shutting down...")
	log.Info().Msg("goodbye!")
}

type verifier struct {
	settingsPermissions rego.PreparedEvalQuery
}

func newVerifier(mainCtx context.Context, regoPath string) (*verifier, error) {
	// --------------------------------------------
	// Set ups
	// --------------------------------------------

	if regoPath == "" {
		return nil, fmt.Errorf("no rego path passed")
	}

	{
		finfo, err := os.Stat(regoPath)
		if err != nil {
			return nil, fmt.Errorf(`could not load "%s": %w`, regoPath, err)
		}

		if !finfo.IsDir() {
			return nil, fmt.Errorf(`"%s" is not a directory`, regoPath)
		}
	}

	bundles := rego.LoadBundle(regoPath)

	// --------------------------------------------
	// Prepare evaluators
	// --------------------------------------------

	settingsPermissions := rego.New(rego.Query("data.users.settings.permissions"), bundles)
	spEval, err := settingsPermissions.PrepareForEval(mainCtx)
	if err != nil {
		return nil, fmt.Errorf(`could not set up "settings.permissions" evaluator: %w`, err)
	}

	return &verifier{
		settingsPermissions: spEval,
	}, nil
}

func (v *verifier) verifySettingsPermissions(ctx context.Context, data []byte) (rego.ResultSet, error) {
	var input interface{}
	if err := util.Unmarshal(data, &input); err != nil {
		return nil, fmt.Errorf("unable to parse input: %w", err)
	}

	inputValue, err := ast.InterfaceToValue(input)
	if err != nil {
		return nil, fmt.Errorf("unable to process input: %w", err)
	}

	return v.settingsPermissions.Eval(ctx, rego.EvalParsedInput(inputValue))
}
