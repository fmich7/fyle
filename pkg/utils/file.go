package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
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

// LocationOnServer checks if the location is on the server is valid
// Returns the correct location if using shorthand
func LocationOnServer(baseDir, user, location, filename string) (string, bool) {
	// Reject unsafe usernames
	user = filepath.Clean(user)
	if user == ".." || strings.Contains(user, "/") || strings.Contains(user, "\\") {
		return "", false //
	}

	// Ensure user directory is set correctly
	userDir, err := filepath.Abs(filepath.Join(baseDir, user))
	if err != nil {
		return "", false
	}

	// Create user directory if it doesn't exist
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		os.Mkdir(userDir, os.ModePerm)
	}

	// Handle location shorthands (`~/`, `.`, `/`) → Defaults to user directory
	homeShorthands := []string{"~", "~/.", ".", "./", "/"}
	for _, shorthand := range homeShorthands {
		if location == shorthand {
			location = ""
			break
		}
	}

	// Join user directory with location and filename
	fullPath := filepath.Join(userDir, location, filename)

	// Normalize path
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", false
	}

	// Ensure the resolved path is still within userDir (prevent directory traversal)
	if !strings.HasPrefix(absPath, userDir) {
		return "", false
	}

	return absPath, true
}
