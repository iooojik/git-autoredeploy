package main

import (
	"flag"

	"git-autoredeploy/internal/checker"
	"git-autoredeploy/internal/config"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	configDir := flag.String("configDir", ".", "Directory where the config file is located")
	flag.Parse()

	cfg := config.LoadConfig(*configDir)

	c := checker.New(cfg)

	viper.OnConfigChange(func(e fsnotify.Event) {
		if e.Op == fsnotify.Remove {
			return
		}

		_ = viper.Unmarshal(cfg)
	})

	viper.WatchConfig()

	err := c.RunMonitoring()
	if err != nil {
		panic(err)
	}
}
