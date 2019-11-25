package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const permDir = 0700
const permFile = 0644

const nilConfig = `{
  "profiles": []
}`

// The Config struct represents the configuration file that tfdd uses to get
// AWS profile information.
type Config struct {
	Path     string     `json:"-"`
	Profiles []*Profile `json:"profiles"`
}

// The Profile struct represents a single AWS profile as a name and AWS account
// ID.
type Profile struct {
	Name      string `json:"name"`
	AccountID string `json:"account_id"`
}

// New creates a new Config object. It checks for a file stored at the filepath
// `cfgPath`, and if it doesn't exist, it creates a new one.
func New(cfgPath string) (*Config, error) {
	var cfg *Config

	_, err := os.Stat(cfgPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to check config file %v", cfgPath))
	}

	if os.IsNotExist(err) {
		folderPath := filepath.Dir(cfgPath)
		err = os.MkdirAll(folderPath, permDir)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to create folder %v", folderPath))
		}

		err = ioutil.WriteFile(cfgPath, []byte(nilConfig), permFile)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to create file %v", cfgPath))
		}
	}

	configFile, err := os.Open(cfgPath) //nolint:gosec
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to open file %v", cfgPath))
	}
	defer configFile.Close()

	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&cfg)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to decode file %v", cfgPath))
	}

	cfg.Path = cfgPath
	return cfg, nil
}

// Save writes a config object to its designated filepath.
func (cfg *Config) Save() error {
	b, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	err = ioutil.WriteFile(cfg.Path, b, permFile)
	return errors.Wrap(err, fmt.Sprintf("failed to write file %v", cfg.Path))
}
