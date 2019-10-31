package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const PERM_DIR = 0700
const PERM_FILE = 0644

const NIL_CONFIG = `{
  "profiles": []
}`

type Config struct {
	Path     string     `json:"-"`
	Profiles []*Profile `json:"profiles"`
}

type Profile struct {
	Name      string `json:"name"`
	AccountID string `json:"account_id"`
}

func New(cfgPath string) (*Config, error) {
	var cfg *Config

	_, err := os.Stat(cfgPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to check config file %v", cfgPath))
	}

	if os.IsNotExist(err) {
		folderPath := filepath.Dir(cfgPath)
		err = os.MkdirAll(folderPath, PERM_DIR)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to create folder %v", folderPath))
		}

		err = ioutil.WriteFile(cfgPath, []byte(NIL_CONFIG), PERM_FILE)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("failed to create file %v", cfgPath))
		}
	}

	configFile, err := os.Open(cfgPath)
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

func (cfg *Config) Save() error {
	b, err := json.MarshalIndent(cfg, "", " ")
	if err != nil {
		return errors.Wrap(err, "failed to marshal config")
	}

	err = ioutil.WriteFile(cfg.Path, b, PERM_FILE)
	return errors.Wrap(err, fmt.Sprintf("failed to write file %v", cfg.Path))
}
