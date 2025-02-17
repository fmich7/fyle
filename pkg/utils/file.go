package utils

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// GetBaseDir returns absolute path without filename
func GetBaseDir(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	fmt.Println(absPath)
	if err != nil {
		return "", err
	}

	info, err := os.Stat(absPath)
	if err != nil {
		return "", err
	}

	if info.IsDir() {
		return absPath, nil
	}

	return filepath.Dir(absPath), nil
}

// Replaces aliases that user might use as a home dir
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

func JoinPathParts(storageRootPath string, parts ...string) string {
	buf := new(bytes.Buffer)
	buf.WriteString(storageRootPath)
	for _, part := range parts {
		// Idk if i should do this together
		buf.WriteString("/")
		buf.WriteString(part)
	}

	joinedPath := filepath.Clean(string(buf.Bytes()))
	return joinedPath
}

// ValidatePath check if constructedPath isn't "higher" than rootPath
func ValidatePath(storageRootPath, constructedPath string) bool {
	// Prevent directory traversal attacks
	if !strings.HasPrefix(constructedPath, storageRootPath) {
		return false
	}

	return true
}

// GetLocationOnServer return joined file location on the server
// Also it checks if the path is not outside user folder
func GetLocationOnServer(storageRootAbsPath, username, subfolders, filename string) (string, bool) {
	// User home dir
	userDir := filepath.Join(storageRootAbsPath, username)

	// Replace aliases if they were provided
	subfolders = ReplaceHomeDirAliases(subfolders)

	fullPath := JoinPathParts(userDir, subfolders, filename)
	valid := ValidatePath(userDir, fullPath)
	if !valid {
		return "", false
	}

	return fullPath, true
}

// GetFileNameFromPath return filename with extension from given path
func GetFileNameFromPath(path string) string {
	return filepath.Base(path)
}

// GetFileNameFromContentDisposition returns filename from Content-Disposition header
func GetFileNameFromContentDisposition(header string) (string, error) {
	args := strings.Split(header, " ")
	target := args[1]
	if target == "" {
		return "", fmt.Errorf("invalid header")
	}
	equalCharIndex := strings.Index(target, "=")
	lastCharIndex := len(target)

	if equalCharIndex == lastCharIndex {
		return "", fmt.Errorf("invalid header")
	}
	return target[equalCharIndex+1 : lastCharIndex], nil
}

// SaveFileOnDisk saves file on disk given its path and content :)
func SaveFileOnDisk(path, filename string, content io.ReadCloser) error {
	newFilePath := JoinPathParts(path, filename)

	// TODO: Dont replace file if exists
	defer content.Close()
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}

	file, err := os.Create(newFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, content)
	if err != nil {
		return err
	}

	return nil
}
