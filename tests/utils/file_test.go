package utils_test

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/fmich7/fyle/pkg/file"
	"github.com/fmich7/fyle/pkg/server"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGetBaseDir(t *testing.T) {
	assert := assert.New(t)

	// test with an existing directory
	tempDir := t.TempDir()
	baseDir, err := server.GetBaseDir(tempDir)
	assert.NoError(err, "expected no error")
	assert.Equal(tempDir, baseDir, "expected base directory to match input directory")

	// test with a file inside the directory
	tempFile := filepath.Join(tempDir, "testfile.txt")
	err = os.WriteFile(tempFile, []byte("test content"), 0644)
	assert.NoError(err, "failed to create temp file")

	baseDir, err = server.GetBaseDir(tempFile)
	assert.NoError(err, "expected no error")
	assert.Equal(tempDir, baseDir, "expected base directory to be parent of the file")

	// test with a non-existent path
	nonExistentPath := filepath.Join(tempDir, "LOL")
	baseDir, err = server.GetBaseDir(nonExistentPath)
	assert.NoError(err, "expected no error for non-existent file")
	assert.Equal(tempDir, baseDir, "expected base directory to be parent of non-existent file")
}

func TestReplaceHomeDirAliases(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		path     string
		expected string
	}{
		{"~", ""},
		{"~/", "/"},
		{"~/.config/", "/.config/"},
		{"./", ""},
		{".", "."},
		{"/", "/"},
		{"home/", "/"},
		{"home/user/", "/user/"},
	}

	for _, test := range tests {
		result := server.ReplaceHomeDirAliases(test.path)
		assert.Equal(test.expected, result, "Path is not equal for input %s", test.path)
	}
}

func TestValidatePath(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		storageRootPath  string
		constructedPath  string
		expectedValidity bool
	}{
		// valid paths
		{"/home/user/storage", "/home/user/storage/file.txt", true},
		{"/home/user/storage", "/home/user/storage/folder", true},
		{"/home/user/storage", "/home/user/storage/folder/file.txt", true},

		// invalid paths (directory traversal)
		{"/home/user/storage", "/home/user/other/file.txt", false},
		{"/home/user/storage", "/home/user/../../etc/passwd", false},

		// edge cases
		{"/home/user/storage", "/home/user/storage/../file.txt", false},
		{"/home/user/storage", "/home/user/storage/./file.txt", true},
	}

	for _, test := range tests {
		result := server.ValidatePath(test.storageRootPath, test.constructedPath)
		assert.Equal(
			test.expectedValidity,
			result,
			"Failed for path: %s with root: %s",
			test.constructedPath,
			test.storageRootPath,
		)
	}
}

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
			// call the function being tested
			safePath, ok := server.GetLocationOnServer(
				test.baseDir,
				test.user,
				test.location,
				test.filename,
			)

			// check if the function returned the expected result
			assert.True(t, ok)
			assert.Equal(t, test.expected, safePath)
		})
	}
}

func TestLocationOnServerUnsafe(t *testing.T) {
	// get the absolute path of the current working directory
	baseDirAbs, err := os.Getwd()
	assert.NoError(t, err, "couldn't get working dir")

	// define the root storage path
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
			_, ok := server.GetLocationOnServer(
				test.rootAbsPath,
				test.username,
				test.subfolders,
				test.filename,
			)
			assert.False(t, ok)
		})
	}
}

func TestGetFileNameFromPath(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		input    string
		expected string
	}{
		{"/path/to/file.txt", "file.txt"},
		{"file.txt", "file.txt"},
		{"/path/to/file", "file"},
		{"/path/.hidden", ".hidden"},
		{".hidden", ".hidden"},
		{"/path/to/", ""},
		{"/", ""},
		{"", ""},
		{"C:/", ""},
		{"C:\\path\\to\\file.txt", "file.txt"},
	}

	for _, test := range tests {
		result := server.GetFileNameFromPath(test.input)
		assert.Equal(test.expected, result, "fail: %s not equal %s", result, test.expected)
	}
}

func TestGetFileNameFromContentDisposition(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		header   string
		expected string
		hasError bool
	}{
		{
			"attachment; filename=example.txt",
			"example.txt",
			false,
		},
		{
			"attachment; filename=example.txt ",
			"example.txt",
			false,
		},
		{
			"attachment; filename=example.txt;",
			"example.txt",
			false,
		},
		{
			"attachment; filename",
			"",
			true,
		},
	}

	for _, test := range tests {
		result, err := server.GetFileNameFromContentDisposition(test.header)

		if test.hasError {
			assert.Error(err, "Expected an error for header %s", test.header)
		} else {
			assert.NoError(err, "Unexpected error for header %s", test.header)
			assert.Equal(test.expected, result, "Filename is not equal for header %s", test.header)
		}
	}
}

func TestSaveFileOnDisk(t *testing.T) {
	assert := assert.New(t)
	afs := afero.NewMemMapFs()
	tempDir := "/temp"
	filename := "testttt.txt"
	path := filepath.Join(tempDir, filename)
	content := []byte("mega content")
	contentReader := io.NopCloser(bytes.NewReader(content))

	// should be created successfully
	err := file.SaveFileOnDisk(afs, path, contentReader)
	assert.NoError(err, "Expected no error when saving file that doesn't exist")
	exists, err := afero.Exists(afs, path)
	assert.NoError(err)
	assert.True(exists, "Expected file to be created")

	// does file content match?
	file, err := afero.ReadFile(afs, path)
	assert.NoError(err, "Failed to read file from afero")
	assert.Equal(content, file, "File content does not match")
}
