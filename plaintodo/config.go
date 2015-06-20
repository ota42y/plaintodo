package main

import (
	"github.com/BurntSushi/toml"
	"io"
)

type Config struct {
	Archive ArchiveConfig
	Paths   PathConfig
	Writer  io.Writer
}

type PathConfig struct {
	Task    string
	History string
}

type ArchiveConfig struct {
	Directory  string
	NameFormat string
}

func readConfig(tomlFilepath string) *Config {
	var config Config
	_, err := toml.DecodeFile(tomlFilepath, &config)
	if err != nil {
		return nil
	}

	return &config
}
