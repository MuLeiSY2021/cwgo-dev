/*
 * Copyright 2022 CloudWeGo Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package config

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/cloudwego/cwgo/pkg/consts"
	"gopkg.in/yaml.v3"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

var globalArgs = NewArgument()

func GetGlobalArgs() *Argument {
	return globalArgs
}

type Argument struct {
	*BasicArgument
	*ServerArgument `yaml:",inline"`
	*ClientArgument
	*ModelArgument
	*DocArgument
	*JobArgument
	*ApiArgument
	*FallbackArgument
}

func NewArgument() *Argument {
	return &Argument{
		BasicArgument:    NewBasicArgument(),
		ServerArgument:   NewServerArgument(),
		ClientArgument:   NewClientArgument(),
		ModelArgument:    NewModelArgument(),
		DocArgument:      NewDocArgument(),
		JobArgument:      NewJobArgument(),
		ApiArgument:      NewApiArgument(),
		FallbackArgument: NewFallbackArgument(),
	}
}

type DialectorFunc func(string) gorm.Dialector

var OpenTypeFuncMap = map[consts.DataBaseType]DialectorFunc{
	consts.MySQL:     mysql.Open,
	consts.SQLServer: sqlserver.Open,
	consts.Sqlite:    sqlite.Open,
	consts.Postgres:  postgres.Open,
}

func WarpArgument(fileExt string, configData []byte) error {
	var err error

	// Call different parsing functions based on the file extension
	switch fileExt {
	case consts.Json:
		globalArgs, err = parseJSON(configData)
	case consts.Yaml, consts.Yml:
		globalArgs, err = parseYAML(configData)
	case consts.Xml:
		globalArgs, err = parseXML(configData)
	default:
		return fmt.Errorf("unsupported file format: %s", fileExt)
	}

	return err
}

// parseJSON parses the JSON configuration data into the Argument struct
func parseJSON(data []byte) (*Argument, error) {
	var args Argument
	err := json.Unmarshal(data, &args)
	if err != nil {
		return nil, fmt.Errorf("failed to parse JSON config: %v", err)
	}
	return &args, nil
}

// parseYAML parses the YAML configuration data into the Argument struct
func parseYAML(data []byte) (*Argument, error) {
	var args Argument
	err := yaml.Unmarshal(data, globalArgs)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %v", err)
	}
	return &args, nil
}

// parseXML parses the XML configuration data into the Argument struct
// Assumes the config.Argument and its sub-structures need to be mapped using XML tags
func parseXML(data []byte) (*Argument, error) {
	var args Argument
	err := xml.Unmarshal(data, &args)
	if err != nil {
		return nil, fmt.Errorf("failed to parse XML config: %v", err)
	}
	return &args, nil
}
