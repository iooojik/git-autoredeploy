package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

type Project struct {
	Name      string `mapstructure:"name"`
	Repo      string `mapstructure:"repo"`
	Directory string `mapstructure:"directory"`
	Command   string `mapstructure:"command"`
}

type Config struct {
	Projects      []Project `mapstructure:"projects"`
	CheckInterval int       `mapstructure:"check_interval"`
}

var config Config

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml") // Specify the config type as YAML
	viper.AddConfigPath(".")    // Look for the config file in the current directory

	if err := viper.ReadInConfig(); err != nil {
		fmt.Printf("Error reading config file, %s\n", err)
		os.Exit(1)
	}

	viper.Unmarshal(&config)
}

func getLastCommitSha(directory string) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = directory
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

func runCommand(directory, command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func monitorProjects() {
	lastCommitShas := make(map[string]string)

	for {
		for _, project := range config.Projects {
			currentSha, err := getLastCommitSha(project.Directory)
			if err != nil {
				fmt.Printf("Error getting last commit for %s: %v\n", project.Name, err)
				continue
			}

			if lastSha, ok := lastCommitShas[project.Name]; ok && lastSha != currentSha {
				fmt.Printf("New update detected for %s! Running command...\n", project.Name)
				err := runCommand(project.Directory, project.Command)
				if err != nil {
					fmt.Printf("Error running command for %s: %v\n", project.Name, err)
				}
			}

			lastCommitShas[project.Name] = currentSha
		}

		time.Sleep(time.Duration(config.CheckInterval) * time.Second)
	}
}

func main() {
	initConfig()

	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
		viper.Unmarshal(&config)
	})
	viper.WatchConfig()

	monitorProjects()
}
