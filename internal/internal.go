package internal

import "errors"

const (
	// sources config file
	SOURCES_CONFIG = "sources.yaml"

	// default TTL in cases not set (1 minute)
	DEFAULT_TTL = 1
)

var (
	// source validation errors
	ERR_NO_DESCRIPTION     = errors.New("No description is set")
	ERR_BAD_URL            = errors.New("url is appears to be invalid")
	ERR_BAD_EXTRACT_CONFIG = errors.New("no extract 'name' or 'value' key set")
	ERR_BAD_SAVE_CONFIG    = errors.New("no save 'path or 'filename' set")
)

// Sources is representation of the parsed YAML sources configuration file
type Sources []struct {
	Source `yaml:"source"`
}

// Extract represents the fields which will be extracted as name & value
type Extract struct {
	Root  string `yaml:"root"`
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

// Save contains the folder and filename for saving
type Save struct {
	Folder   string `yaml:"folder"`
	Filename string `yaml:"filename"`
}

// Source represents an external data source and the config for parsing it
type Source struct {
	Description string  `yaml:"description"`
	Url         string  `yaml:"url"`
	Token       string  `yaml:"token"`
	Extract     Extract `yaml:"extract"`
	Ttl         int     `yaml:"ttl"`
	Save        Save    `yaml:"save"`
}
