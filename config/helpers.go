package config

import (
	"github.com/BurntSushi/toml"
)

func ParseConfig(path string) error {
	if path == "" {
		path = "./config.toml"
	}

	_, err := toml.DecodeFile(path, &Cfg)
	return err
}
