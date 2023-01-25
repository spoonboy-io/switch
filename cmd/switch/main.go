package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/spoonboy-io/switch/internal/process"

	"github.com/spoonboy-io/switch/internal/source"

	"github.com/spoonboy-io/koan"
	"github.com/spoonboy-io/reprise"
	"github.com/spoonboy-io/switch/internal"
)

var (
	version   = "Development build"
	goversion = "Unknown"
)

var logger *koan.Logger
var config internal.Sources

func init() {
	logger = &koan.Logger{}

	// read sources and validate the config
	cfg, err := source.ReadAndParseConfig(internal.SOURCES_CONFIG)
	if err != nil {
		logger.FatalError("Failed to read sources configuration file", err)
	}
	err = source.ValidateConfig(cfg)
	if err != nil {
		logger.FatalError("Failed to validate sources configuration", err)
	}

	config = cfg
}

func shutdown(cancel context.CancelFunc) {
	fmt.Println("") // break after ^C
	logger.Warn("Application terminated")
	cancel()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer shutdown(cancel)

	// write a console banner
	reprise.WriteSimple(&reprise.Banner{
		Name:         "Switch",
		Description:  "Preprocess JSON into name and value fields",
		Version:      version,
		GoVersion:    goversion,
		WebsiteURL:   "https://spoonboy.io",
		VcsURL:       "https://github.com/spoonboy-io/switch",
		VcsName:      "Github",
		EmailAddress: "hello@spoonboy.io",
	})

	// refresh loop
	go func() {
		checkInterval := time.NewTicker(61 * time.Second)
		logger.Info("Checking TTLs for stale data")
		process.CheckAndRefresh(ctx, config, logger)
		for range checkInterval.C {
			logger.Info("Checking TTLs for stale data")
			process.CheckAndRefresh(ctx, config, logger)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
