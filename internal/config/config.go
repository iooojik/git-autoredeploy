package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type Project struct {
	Name      string `yaml:"name"`
	Repo      string `yaml:"repo"`
	Directory string `yaml:"directory"`
	Command   string `yaml:"command"`
}

type Config struct {
	Projects      []Project `yaml:"projects"`
	CheckInterval int       `yaml:"check_interval"`
}

func LoadConfig(path string) *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	f, err := os.Open(path)
	if err != nil {
		panic(fmt.Errorf("open config file: %w", err))
	}

	err = viper.ReadConfig(f)
	if err != nil {
		panic(fmt.Errorf("read config file %s %w", path, err))
	}

	cfg := new(Config)

	err = viper.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
