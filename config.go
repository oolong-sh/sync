package sync

import (
	"os"

	"github.com/BurntSushi/toml"
)

type SyncConfig struct {
	Host           string `toml:"host"`
	User           string `toml:"user"`
	Port           int    `toml:"port"`
	PrivateKeyPath string `toml:"private_key_path"`
	// TODO: ssh key password
}

func ReadConfig(configPath string) (SyncConfig, error) {
	cfg := SyncConfig{}

	contents, err := os.ReadFile(configPath)
	if err != nil {
		return cfg, err
	}

	err = toml.Unmarshal(contents, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, err
}
