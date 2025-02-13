package utils

import (
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
			safePath, ok := LocationOnServer(test.baseDir, test.user, test.location, test.filename)

			// Check if the function returned the expected result
			assert.True(t, ok)
			assert.Equal(t, test.expected, safePath)
		})
	}
}

func TestLocationOnServerUnsafe(t *testing.T) {
	tests := []struct {
		testName string
		baseDir  string
		user     string
		location string
		filename string
	}{
		{
			"Location contains ..",
			"/server/uploads",
			"testuser",
			"../documents",
			"file.txt",
		},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			// Call the function being tested
			_, ok := LocationOnServer(test.baseDir, test.user, test.location, test.filename)
			assert.False(t, ok)
		})
	}
}
