package main

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Paths PathConfig
}

type PathConfig struct {
	Task string
}

func readConfig(tomlFilepath string) *Config{
	var config Config
	_, err := toml.DecodeFile(tomlFilepath, &config)
	if err != nil {
		panic(err)
	}

	return &config
}

func main() {
}
