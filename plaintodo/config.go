package main

import (
	"github.com/BurntSushi/toml"
	"io"
)

type Config struct {
	Paths  PathConfig
	Writer io.Writer
}

type PathConfig struct {
	Task string
}

func readConfig(tomlFilepath string) *Config {
	var config Config
	_, err := toml.DecodeFile(tomlFilepath, &config)
	if err != nil {
		return nil
	}

	return &config
}
