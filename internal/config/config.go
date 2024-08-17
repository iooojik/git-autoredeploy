package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

// Project represents a single project configuration with its repository details and command to execute.
type Project struct {
	Name      string `yaml:"name"`
	Repo      string `yaml:"repo"`
	Directory string `yaml:"directory"`
	Command   string `yaml:"command"`
}

// Config holds the configuration for all projects and the interval for checking updates.
type Config struct {
	Projects      []Project `yaml:"projects"`
	CheckInterval int       `yaml:"check_interval"`
}

// LoadConfig loads the configuration from the specified path using Viper.
func LoadConfig(path string) *Config {
	// Set up Viper to look for a config.yaml file in the configs directory
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	// Open the configuration file
	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("open config file: %w", err))
	}

	// Read the configuration file into Viper
	err = viper.ReadConfig(f)
	if err != nil {
		panic(fmt.Errorf("read config file %s %w", path, err))
	}

	// Unmarshal the configuration into the Config struct
	cfg := new(Config)

	err = viper.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
