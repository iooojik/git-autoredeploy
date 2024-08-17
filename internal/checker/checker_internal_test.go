package checker

import (
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestFetchLocalCommitSha tests the fetchLocalCommitSha function.
func TestFetchLocalCommitSha(t *testing.T) {
	t.Parallel()

	t.Run("Valid SHA Retrieval", func(t *testing.T) {
		t.Parallel()

		// Set up a temporary directory and initialize a Git repository
		dir, err := os.MkdirTemp("", "gitrepo")
		require.NoError(t, err)

		defer func() { _ = os.RemoveAll(dir) }()

		cmd := exec.Command("git", "init", dir)
		err = cmd.Run()
		require.NoError(t, err)

		// Create a dummy commit
		//nolint:revive
		cmd = exec.Command("sh", "-c", "echo 'test' > testfile && git add testfile && git commit -m 'Initial commit'")

		cmd.Dir = dir
		err = cmd.Run()
		require.NoError(t, err)

		sha, err := fetchLocalCommitSha(dir)
		require.NoError(t, err)
		require.Len(t, sha, 40)
	})
}

// TestRunCommand tests the runCommand function.
func TestRunCommand(t *testing.T) {
	t.Parallel()

	t.Run("Valid Command Execution", func(t *testing.T) {
		t.Parallel()

		// Set up a temporary directory
		dir, err := os.MkdirTemp("", "testcmd")
		require.NoError(t, err)

		defer func() { _ = os.RemoveAll(dir) }()

		// Test a simple echo command
		err = runCommand(dir, "echo 'hello world'")
		require.NoError(t, err)
	})
}
