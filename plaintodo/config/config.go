package config

import (
	"io"
	"path"

	"github.com/BurntSushi/toml"
)

// Config is config data struct
type Config struct {
	Archive ArchiveConfig
	Task    TaskConfig
	Liner   LinerConfig
	Writer  io.Writer
	Command CommandConfig
}

// LinerConfig is liner^s config
type LinerConfig struct {
	History string
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

// TaskConfig is task setting file
type TaskConfig struct {
	TaskFolder      string
	DefaultFilename string
}

// GetDefaultTaskFilepath return default task filepath set in config.toml
func (c TaskConfig) GetDefaultTaskFilepath() string {
	return c.GetTaskFilepath(c.DefaultFilename)
}

// GetTaskFilepath return selected task filepath in task folder
func (c TaskConfig) GetTaskFilepath(filename string) string {
	return path.Join(c.TaskFolder, filename)
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
