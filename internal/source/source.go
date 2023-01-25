package source

import (
	"io/ioutil"
	"net/url"

	"github.com/spoonboy-io/switch/internal"

	"gopkg.in/yaml.v2"
)

var config internal.Sources

// ReadAndParseConfig reads the contents of the YAML source config file
// and parses it to a map of Source structs
func ReadAndParseConfig(cfgFile string) error {
	yamlConfig, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		return err
	}

	if err := yaml.Unmarshal(yamlConfig, &config); err != nil {
		return err
	}

	return nil
}

func ValidateConfig() error {
	for i := range config {
		// check description
		if config[i].Description == "" {
			return internal.ERR_NO_DESCRIPTION
		}

		// check method
		if err := isGoodMethod(config[i].Method); err != nil {
			return err
		}

		// check url
		if _, err := url.ParseRequestURI(config[i].Url); err != nil {
			return internal.ERR_BAD_URL
		}

		// if method POST/PUT check request body is present
		if err := shouldHaveRequestBody(config[i].Method, config[i].RequestBody); err != nil {
			return err
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

// helpers
func isGoodMethod(method string) error {
	switch method {
	case "GET", "POST":
		return nil
	default:
		return internal.ERR_BAD_METHOD
	}
}

func shouldHaveRequestBody(method, requestBody string) error {
	if method != "GET" {
		if requestBody == "" {
			return internal.ERR_NO_BODY
		}
	}
	return nil
}
