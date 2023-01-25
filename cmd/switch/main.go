package main

import (
	"os"
	"path/filepath"

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

func init() {
	logger = &koan.Logger{}

	// cache folder for sources which we set a TTL for
	dataPath := filepath.Join(".", internal.CACHE_FOLDER)
	if err := os.MkdirAll(dataPath, os.ModePerm); err != nil {
		logger.FatalError("Problem checking/creating 'cache' folder", err)
	}

	// read sources and validate the config
	err := source.ReadAndParseConfig(internal.SOURCES_CONFIG)
	if err != nil {
		logger.FatalError("Failed to read sources configuration file", err)
	}
	err = source.ValidateConfig()
	if err != nil {
		logger.FatalError("Failed to validate sources configuration", err)
	}
}

func main() {
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
}
