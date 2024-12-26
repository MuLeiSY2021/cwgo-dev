package arg_config

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/cloudwego/cwgo/config"
	"github.com/cloudwego/cwgo/pkg/consts"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

// ReadConfig reads the configuration file based on its extension and parses it accordingly
func ReadConfig() (*config.Argument, error) {
	path := config.BasicArguments.Config
	// Determine the file extension to decide which parser to use
	fileExt := strings.ToLower(path[strings.LastIndex(path, ".")+1:])
	var configData []byte
	var err error

	// Read the file
	if configData, err = os.ReadFile(path); err != nil {
		return nil, fmt.Errorf("failed to read the config file: %v", err)
	}

	// Call different parsing functions based on the file extension
	var arguments *config.Argument

	switch fileExt {
	case consts.Json:
		arguments, err = parseJSON(configData)
	case consts.Yaml, consts.Yml:
		arguments, err = parseYAML(configData)
	case consts.Xml:
		arguments, err = parseXML(configData)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", fileExt)
	}

	if err != nil {
		return nil, err
	}

	return arguments, nil
}

// parseJSON parses the JSON configuration data into the Argument struct
func parseJSON(data []byte) (*config.Argument, error) {
	var args config.Argument
	err := json.Unmarshal(data, &args)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON config: %v", err)
	}
	return &args, nil
}

// parseYAML parses the YAML configuration data into the Argument struct
func parseYAML(data []byte) (*config.Argument, error) {
	var args config.Argument
	err := yaml.Unmarshal(data, &args)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %v", err)
	}
	return &args, nil
}

// parseXML parses the XML configuration data into the Argument struct
// Assumes the config.Argument and its sub-structures need to be mapped using XML tags
func parseXML(data []byte) (*config.Argument, error) {
	var args config.Argument
	err := xml.Unmarshal(data, &args)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML config: %v", err)
	}
	return &args, nil
}
