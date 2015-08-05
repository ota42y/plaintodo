package config

import (
	"github.com/BurntSushi/toml"
	"io"
)

// Config is config data struct
type Config struct {
	Archive ArchiveConfig
	Paths   PathConfig
	Writer  io.Writer
	Command CommandConfig
}

// Alias is command alias data in config file
type Alias struct {
	Name    string
	Command string
}

// CommandConfig is Alias command config in config file
type CommandConfig struct {
	Aliases []Alias
}

// PathConfig is path setting in config file
type PathConfig struct {
	Task    string
	History string
}

// ArchiveConfig is Archive config in config file
type ArchiveConfig struct {
	Directory  string
	NameFormat string
}

// ReadConfig read config from file
func ReadConfig(tomlFilepath string) *Config {
	var config Config
	_, err := toml.DecodeFile(tomlFilepath, &config)
	if err != nil {
		return nil
	}

	return &config
}
