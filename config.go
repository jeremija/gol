package gol

import (
	"github.com/BurntSushi/toml"
	"github.com/jeremija/gol/dispatchers"
)

type AppConfig struct {
	Files      []*FileTailerConfig
	Dispatcher dispatchers.DispatcherConfig
	DryRun     bool
}

func ReadConfig(file string) (AppConfig, error) {
	if file == "" {
		return AppConfig{
			Files: make([]*FileTailerConfig, 0),
		}, nil
	}
	var config AppConfig
	_, err := toml.DecodeFile(file, &config)
	return config, err
}
