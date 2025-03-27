package utils

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// GetBaseDir returns absolute path without filename (parent folder).
func GetBaseDir(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		// if the path doesn't exist, return its parent directory
		if os.IsNotExist(err) {
			return filepath.Dir(absPath), nil
		}
		return "", err
	}

	if info.IsDir() {
		return absPath, nil
	}

	return filepath.Dir(absPath), nil
}

// Replaces aliases that user might use as a home dir.
func ReplaceHomeDirAliases(path string) string {
	// handle home shortcuts -> defaults to user directory
	homeShortcuts := []string{"~", "home", "./"}
	for _, shorthand := range homeShortcuts {
		path = strings.TrimPrefix(path, shorthand)
	}

	return path
}

// JoinPathParts joins path parts with starting path.
func JoinPathParts(storageRootPath string, parts ...string) string {
	// join parts
	joinedParts := filepath.Join(parts...)

	// join storage root path with joined parts
	finalPath := filepath.Join(storageRootPath, joinedParts)

	return filepath.Clean(finalPath)
}

// ValidatePath check if constructedPath isn't "higher" than rootPath.
func ValidatePath(storageRootPath, constructedPath string) bool {
	path := filepath.Clean(constructedPath)

	// prevent directory traversal attacks
	if !strings.HasPrefix(path, storageRootPath) {
		return false
	}

	return true
}

// GetLocationOnServer return joined file location on the server.
// Also it checks if the path is not outside user folder.
func GetLocationOnServer(storageRootAbsPath, username, subfolders, filename string) (string, bool) {
	// user home dir
	userDir := filepath.Join(storageRootAbsPath, username)

	// replace aliases if they were provided
	subfolders = ReplaceHomeDirAliases(subfolders)

	fullPath := JoinPathParts(userDir, subfolders, filename)
	valid := ValidatePath(userDir, fullPath)
	if !valid {
		return "", false
	}

	return fullPath, true
}

// GetFileNameFromPath return filename with extension from given path.
func GetFileNameFromPath(path string) string {
	// windows paths to Unix style
	path = strings.ReplaceAll(path, "\\", "/")

	// is dir
	if path == "" || path[len(path)-1] == os.PathSeparator {
		return ""
	}

	return filepath.Base(path)
}

// GetFileNameFromContentDisposition returns filename from Content-Disposition header.
func GetFileNameFromContentDisposition(header string) (string, error) {
	lowerHeader := strings.ToLower(header)
	if idx := strings.Index(lowerHeader, "filename="); idx != -1 {
		start := idx + len("filename=")
		filename := header[start:]

		// ; after filename
		if idx = strings.Index(filename, ";"); idx != -1 {
			filename = filename[:idx]
		}

		// " " space after filename
		if idx = strings.Index(filename, " "); idx != -1 {
			filename = filename[:idx]
		}

		return strings.TrimSpace(filename), nil
	}

	return "", errors.New("invalid header")
}
