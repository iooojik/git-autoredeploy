package main

import (
	"git-autoredeploy/internal/checker"
	"git-autoredeploy/internal/config"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	cfg := config.LoadConfig()

	c := checker.New(cfg)

	viper.OnConfigChange(func(e fsnotify.Event) {
		_ = viper.Unmarshal(cfg)
	})

	viper.WatchConfig()

	err := c.RunMonitoring()
	if err != nil {
		panic(err)
	}
}
