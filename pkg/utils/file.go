package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// GetAbsolutePath returns the absolute path of a file
func GetAbsolutePath(path string) (string, error) {
	// Get working directory
	wd, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("Error getting working directory: %s", err)
	}

	// Join path
	absPath, err := filepath.Abs(filepath.Join(wd, path))
	if err != nil {
		return "", fmt.Errorf("Error getting absolute path: %s", err)
	}

	return absPath, nil
}
