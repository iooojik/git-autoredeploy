package main

import (
	"flag"

	"git-autoredeploy/internal/checker"
	"git-autoredeploy/internal/config"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

func main() {
	// Define a command-line flag for the config directory, defaulting to the current directory
	configDir := flag.String("config", "configs/config.yaml", "Directory where the config file is located")
	flag.Parse()

	// Load the configuration from the specified directory
	cfg := config.LoadConfig(*configDir)

	// Create a new checker instance with the loaded configuration
	c := checker.New(cfg)

	// Set up a watcher to monitor changes to the config file
	viper.OnConfigChange(func(e fsnotify.Event) {
		// Ignore the event if the config file was removed
		if e.Op == fsnotify.Remove {
			return
		}

		// Unmarshal the updated config into the cfg struct
		_ = viper.Unmarshal(cfg)
	})

	// Start watching the config file for changes
	viper.WatchConfig()

	// Run the monitoring process defined in the checker
	err := c.RunMonitoring()
	if err != nil {
		panic(err)
	}
}
