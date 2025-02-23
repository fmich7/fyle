package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocationOnServer(t *testing.T) {
	tests := []struct {
		baseDir  string
		user     string
		location string
		filename string
		expected string
	}{
		{
			"/server/uploads",
			"testuser",
			"documents",
			"file.txt",
			"/server/uploads/testuser/documents/file.txt",
		},
		{
			"/server/uploads",
			"testuser",
			".",
			"file.txt",
			"/server/uploads/testuser/file.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			// Call the function being tested
			safePath, ok := GetLocationOnServer(test.baseDir, test.user, test.location, test.filename)

			// Check if the function returned the expected result
			assert.True(t, ok)
			assert.Equal(t, test.expected, safePath)
		})
	}
}

func TestLocationOnServerUnsafe(t *testing.T) {
	// Get the absolute path of the current working directory
	baseDirAbs, err := os.Getwd()
	assert.NoError(t, err, "error getting working dir")

	// Define the root storage path
	storageRootAbsPath := filepath.Join(baseDirAbs, "uploads")

	tests := []struct {
		testName    string
		rootAbsPath string
		username    string
		subfolders  string
		filename    string
	}{
		{
			"Location contains ..",
			storageRootAbsPath,
			"testuser",
			"../documents",
			"file.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Call the function being tested
			_, ok := GetLocationOnServer(test.rootAbsPath, test.username, test.subfolders, test.filename)
			assert.False(t, ok)
		})
	}
}
