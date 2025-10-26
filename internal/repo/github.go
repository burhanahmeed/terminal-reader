package repo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func CloneGithub(repoURL *string) (*string, error) {
	tmpDir := filepath.Join(os.TempDir(), "ragrepo_terminal_reader")
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		if err := os.MkdirAll(tmpDir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create temp directory: %w", err)
		}
	}

	repoName := filepath.Base(*repoURL)
	if repoName[len(repoName)-4:] == ".git" {
		repoName = repoName[:len(repoName)-4]
	}
	destPath := filepath.Join(tmpDir, repoName)

	// If already cloned, pull latest
	if _, err := os.Stat(destPath); os.IsNotExist(err) {
		fmt.Println("Cloning GitHub repo:", *repoURL)
		cmd := exec.Command("git", "clone", "--depth", "1", *repoURL, destPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return nil, err
		}
	} else {
		fmt.Println("Repo already cloned, pulling latest")
		cmd := exec.Command("git", "-C", destPath, "pull")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Run() // ignore pull error for now
	}

	return &destPath, nil
}
