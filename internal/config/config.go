package config

import (
	"fmt"

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

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("configs")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err))
	}

	cfg := new(Config)

	err = viper.Unmarshal(cfg)
	if err != nil {
		panic(err)
	}

	return cfg
}
