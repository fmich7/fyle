package utils

import (
	"bytes"
	"path/filepath"
	"strings"
)

// ReplaceHomeDirAliases replaces aliases that user might use as a home dir
func ReplaceHomeDirAliases(path string) string {
	res := path
	// Handle home shortcuts -> defaults to user directory
	homeShortcuts := []string{"~", "~/.", ".", "./", "/", "home"}
	for _, shorthand := range homeShortcuts {
		if res == shorthand {
			return ""
		}
	}

	return res
}

// Checks if constructed path corresponds to a file stored in storageRootPath
func validateAndJoinParts(storageRootPath string, parts ...string) (string, bool) {
	buf := &bytes.Buffer{}
	buf.WriteString(storageRootPath)
	for _, part := range parts {
		// Idk if i should do this together
		buf.WriteString("/")
		buf.WriteString(part)
	}

	joinedPath := filepath.Clean(string(buf.Bytes()))

	// Prevent directory traversal attacks
	if !strings.HasPrefix(joinedPath, storageRootPath) {
		return "", false
	}

	return joinedPath, true
}

// GetLocationOnServer return joined file location on the server
// Also it checks if the path is not outside user folder
func GetLocationOnServer(storageRootAbsPath, username, subfolders, filename string) (string, bool) {
	// User home dir
	userDir := filepath.Join(storageRootAbsPath, username)

	// Replace aliases if they were provided
	subfolders = ReplaceHomeDirAliases(subfolders)

	fullPath, valid := validateAndJoinParts(userDir, subfolders, filename)
	if !valid {
		return "", false
	}

	return fullPath, true
}

// GetFileNameFromPath return filename with extension from given path
func GetFileNameFromPath(path string) string {
	return filepath.Base(path)
}
