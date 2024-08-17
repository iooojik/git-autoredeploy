package checker

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"git-autoredeploy/internal/config"
)

type Checker struct {
	cfg *config.Config
}

func New(cfg *config.Config) *Checker {
	return &Checker{cfg: cfg}
}

func (c *Checker) RunMonitoring() error {
	lastCommitSha := make(map[string]string)

	for _, project := range c.cfg.Projects {
		currentSha, err := fetchLocalCommitSha(project.Directory)
		if err != nil {
			return fmt.Errorf("get local commit sha for %s: %w", project.Name, err)
		}

		lastCommitSha[project.Name] = currentSha
	}

	ticker := time.NewTicker(time.Duration(c.cfg.CheckInterval) * time.Second)
	defer ticker.Stop()

	for t := range ticker.C {
		log.Printf("check at: %v\n", t.UTC())

		for _, project := range c.cfg.Projects {
			newSha, err := checkProject(lastCommitSha[project.Name], project)
			if err != nil {
				return err
			}

			lastCommitSha[project.Name] = newSha
		}
	}

	return nil
}

func checkProject(
	lastSha string,
	project config.Project,
) (string, error) {
	remoteSha, err := fetchRemoteShaForCurrentBranch(project.Directory)
	if err != nil {
		return "", fmt.Errorf("get last commi for %s: %w", project.Name, err)
	}

	if lastSha != remoteSha {
		err = runCommand(project.Directory, project.Command)
		if err != nil {
			return "", fmt.Errorf("run command for %s: %w", project.Name, err)
		}
	}

	return remoteSha, nil
}

func fetchLocalCommitSha(
	directory string,
) (string, error) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = directory

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("rev-parse command: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

func fetchRemoteShaForCurrentBranch(directory string) (string, error) {
	// Get the current branch name
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchCmd.Dir = directory

	branchOutput, err := branchCmd.Output()
	if err != nil {
		return "", fmt.Errorf("get current branch: %w", err)
	}

	branch := strings.TrimSpace(string(branchOutput))

	// Get the SHA for the current branch from the remote
	shaCmd := exec.Command("git", "ls-remote", "origin", branch)
	shaCmd.Dir = directory

	shaOutput, err := shaCmd.Output()
	if err != nil {
		return "", fmt.Errorf("get remote sha for %s: %w", directory, err)
	}

	// Extract and return the SHA
	sha := strings.Fields(string(shaOutput))[0]

	return sha, nil
}

func runCommand(directory, command string) error {
	cmd := exec.Command("sh", "-c", command)
	cmd.Dir = directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("run: %w", err)
	}

	return nil
}
