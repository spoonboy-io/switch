package source

import (
	"io/ioutil"
	"net/url"

	"github.com/spoonboy-io/switch/internal"

	"gopkg.in/yaml.v2"
)

// ReadAndParseConfig reads the contents of the YAML source config file
// and parses it to a map of Source structs
func ReadAndParseConfig(cfgFile string) (internal.Sources, error) {
	var config internal.Sources
	yamlConfig, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return nil, err
	}

	if err := yaml.Unmarshal(yamlConfig, &config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfig checks the basics are complete for each source in the YAML
// and errors if not, also sets the default Ttl if not specified
func ValidateConfig(config internal.Sources) error {
	for i := range config {
		// check description
		if config[i].Description == "" {
			return internal.ERR_NO_DESCRIPTION
		}

		// check url
		if _, err := url.ParseRequestURI(config[i].Url); err != nil {
			return internal.ERR_BAD_URL
		}

		// check extract name and  value
		if config[i].Extract.Name == "" || config[i].Extract.Value == "" {
			return internal.ERR_BAD_EXTRACT_CONFIG
		}

		// check ttl or set default
		if config[i].Ttl == 0 {
			config[i].Ttl = internal.DEFAULT_TTL
		}

		// check save details
		if config[i].Save.Folder == "" || config[i].Save.Filename == "" {
			return internal.ERR_BAD_SAVE_CONFIG
		}
	}

	return nil
}
