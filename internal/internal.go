package internal

const (
	// caching configuration
	CACHE_FOLDER = "cache"

	// sources config file
	SOURCES_CONFIG = "sources.yaml"
)

// Sources is representation of the parsed YAML sources configuration file
type Sources []struct {
	Source `yaml:"source"`
}

// Extract represents the fields which will be extracted as name & value
type Extract struct {
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
	Method      string  `yaml:"method"`
	RequestBody string  `yaml:"requestBody"`
	Token       string  `yaml:"token"`
	Extract     Extract `yaml:"extract"`
	Ttl         int     `yaml:"ttl"`
	Save        Save    `yaml:"save"`
}
