package file

import (
	"bytes"
	"io"
	"path/filepath"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestSaveFileOnDisk(t *testing.T) {
	assert := assert.New(t)
	afs := afero.NewMemMapFs()
	tempDir := "/temp"
	filename := "testttt.txt"
	path := filepath.Join(tempDir, filename)
	content := []byte("mega content")
	contentReader := io.NopCloser(bytes.NewReader(content))

	// should be created successfully
	err := SaveFileOnDisk(afs, path, contentReader)
	assert.NoError(err, "Expected no error when saving file that doesn't exist")
	exists, err := afero.Exists(afs, path)
	assert.NoError(err)
	assert.True(exists, "Expected file to be created")

	// does file content match?
	file, err := afero.ReadFile(afs, path)
	assert.NoError(err, "Failed to read file from afero")
	assert.Equal(content, file, "File content does not match")
}
