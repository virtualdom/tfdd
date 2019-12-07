package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/virtualdom/tfdd/pkg/config"
)

func getConfig(configFilePath string) (*config.Config, error) {
	if len(configFilePath) == 0 {
		configFilePath = path.Join(os.Getenv("HOME"), ".tfdd", "config")
	}

	configFilePath, err := filepath.Abs(configFilePath)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get absolute path of config %v", configFilePath))
	}

	cfg, err := config.New(configFilePath)
	return cfg, errors.Wrap(err, fmt.Sprintf("failed to load config file %v", configFilePath))
}
