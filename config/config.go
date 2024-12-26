package config

import (
	"fmt"
	"github.com/cloudwego/cwgo/pkg/common/utils"
	"github.com/cloudwego/cwgo/pkg/consts"
	"github.com/cloudwego/cwgo/tpl"
	"os"
	"path"
	"strings"
)

type ConfigArgument struct {
	Config       string
	ConfigName   string
	ConfigBranch string
}

func NewConfigArgument() *ConfigArgument {
	return &ConfigArgument{}
}

// ReadConfig reads the configuration file based on its extension and parses it accordingly
func ReadConfig() error {
	configPath := GetGlobalArgs().Config
	// Determine the file extension to decide which parser to use
	fileExt := strings.ToLower(configPath[strings.LastIndex(configPath, ".")+1:])
	var configData []byte
	var err error

	// Read the file

	if strings.HasSuffix(configPath, consts.SuffixGit) {
		err = utils.GitClone(configPath, path.Join(tpl.HertzDir, consts.Client))
		if err != nil {
			return err
		}
		gitPath, err := utils.GitPath(configPath)
		if err != nil {
			return err
		}
		gitPath = path.Join(tpl.ConfigDir, gitPath)
		if err = utils.GitCheckout(GetGlobalArgs().ConfigBranch, gitPath); err != nil {
			return err
		}
		configPath = path.Join(gitPath, GetGlobalArgs().ConfigName)
	}

	if configData, err = os.ReadFile(configPath); err != nil {
		return fmt.Errorf("failed to read the config file: %v", err)
	}

	err = WarpArgument(fileExt, configData)

	return err
}
