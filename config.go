package gol

import (
	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Files map[string]*FileTailerConfig
}

func ReadConfig(file string) (tomlConfig, error) {
	if file == "" {
		return tomlConfig{
			Files: make(map[string]*FileTailerConfig),
		}, nil
	}
	var config tomlConfig
	_, err := toml.DecodeFile(file, &config)
	return config, err
}
