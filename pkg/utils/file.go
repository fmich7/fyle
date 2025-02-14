package utils

import (
	"path/filepath"
	"strings"
)

// GetLocationOnServer return joined file location on the server
// Also it checks if the path is not outside user folder
func GetLocationOnServer(storageRootAbsPath, username, subfolders, filename string) (string, bool) {
	// User home dir
	userDir := filepath.Join(storageRootAbsPath, username)

	// Handle home shortcuts -> defaults to user directory
	homeShortcuts := []string{"~", "~/.", ".", "./", "/"}
	for _, shorthand := range homeShortcuts {
		if subfolders == shorthand {
			subfolders = ""
			break
		}
	}

	// Join user directory with subfolders and filename
	fullPath := filepath.Join(userDir, subfolders, filename)

	// Normalize path
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", false
	}

	// Prevent directory traversal attacks
	if !strings.HasPrefix(absPath, userDir) {
		return "", false
	}

	return absPath, true
}
