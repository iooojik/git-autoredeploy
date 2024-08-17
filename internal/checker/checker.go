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

// Checker manages the monitoring of Git repositories and the execution of commands.
type Checker struct {
	cfg *config.Config
}

// New creates a new Checker instance with the provided configuration.
func New(cfg *config.Config) *Checker {
	return &Checker{cfg: cfg}
}

// RunMonitoring starts the monitoring process that periodically checks repositories
// and runs commands if new commits are detected.
func (c *Checker) RunMonitoring() error {
	lastCommitSha := make(map[string]string)

	// Initialize the map with the latest commit SHA for each project
	for _, project := range c.cfg.Projects {
		currentSha, err := fetchLocalCommitSha(project.Directory)
		if err != nil {
			return fmt.Errorf("get local commit sha for %s: %w", project.Name, err)
		}

		lastCommitSha[project.Name] = currentSha
	}

	// Set up a ticker to trigger checks at regular intervals
	ticker := time.NewTicker(time.Duration(c.cfg.CheckInterval) * time.Second)
	defer ticker.Stop()

	// Monitor the repositories at each tick
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

// checkProject compares the last known commit SHA with the remote SHA for the project's branch.
// If there is a difference, it runs the associated command.
func checkProject(
	lastSha string,
	project config.Project,
) (string, error) {
	remoteSha, err := fetchRemoteShaForCurrentBranch(project.Directory)
	if err != nil {
		return "", fmt.Errorf("get last commit for %s: %w", project.Name, err)
	}

	// If the remote SHA differs, run the project's command
	if lastSha != remoteSha {
		err = runCommand(project.Directory, project.Command)
		if err != nil {
			return "", fmt.Errorf("run command for %s: %w", project.Name, err)
		}
	}

	return remoteSha, nil
}

// fetchLocalCommitSha retrieves the latest commit SHA for the current branch in the local repository.
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

// fetchRemoteShaForCurrentBranch retrieves the latest commit SHA for the current branch
// from the remote repository.
func fetchRemoteShaForCurrentBranch(directory string) (string, error) {
	// Get the current branch name
	branchCmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	branchCmd.Dir = directory

	branchOutput, err := branchCmd.Output()
	if err != nil {
		return "", fmt.Errorf("get current branch: %w", err)
	}

	branch := strings.TrimSpace(string(branchOutput))

	// Get the latest commit SHA for the current branch from the remote repository
	shaCmd := exec.Command("git", "ls-remote", "origin", branch)
	shaCmd.Dir = directory

	shaOutput, err := shaCmd.Output()
	if err != nil {
		return "", fmt.Errorf("ls-remote command: %w", err)
	}

	// Extract the SHA from the output
	sha := strings.Fields(string(shaOutput))[0]

	return sha, nil
}

// runCommand executes the specified command in the given directory.
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
