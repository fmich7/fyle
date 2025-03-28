package storage

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fmich7/fyle/pkg/file"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDiskStorage(t *testing.T) {
	// afero
	afs := afero.NewMemMapFs()
	_, err := NewDiskFileStorage("uploads", afs)
	require.NoError(t, err, "failed to create storage: %v", err)
	exists, err := afero.Exists(afs, "uploads")
	require.NoError(t, err)
	assert.True(t, exists)

	// simulate real os
	wd, err := os.Getwd()
	require.NoError(t, err, "failed to get working directory: %v", err)

	path := filepath.Join(wd, "uploads")
	_, err = NewDiskFileStorage(path, afs)
	require.NoError(t, err, "failed to create storage: %v", err)

	exists, err = afero.Exists(afs, path)
	require.NoError(t, err)
	assert.True(t, exists)
}

func TestStoreFile(t *testing.T) {
	afs := afero.NewMemMapFs()
	storage, err := NewDiskFileStorage("uploads", afs)
	require.NoError(t, err, "failed to create storage: %v", err)

	filename := "test.txt"
	content := []byte("some content")

	err = afero.WriteFile(afs, filename, content, 0777)
	require.NoError(t, err, "failed to create test file: %v", err)

	oFile, err := afs.Open(filename)
	require.NoError(t, err, "failed to open test file: %v", err)

	targetFilename := "TESTME.txt"
	f := file.File{
		Location: targetFilename,
		Data:     oFile,
	}

	err = storage.StoreFile(&f)
	require.NoError(t, err, "failed to save file: %v", err)

	exists, err := afero.Exists(afs, targetFilename)
	require.NoError(t, err, "failed to check file existence: %v", err)
	require.True(t, exists, "expected file to exist after storing")

	savedContent, err := afero.ReadFile(afs, targetFilename)
	require.NoError(t, err, "failed to read stored file: %v", err)

	assert.Equal(t, content, savedContent, "stored content does not match original content")
}

func TestRetrieveFile(t *testing.T) {
	afs := afero.NewMemMapFs()
	storage, err := NewDiskFileStorage("uploads", afs)
	require.NoError(t, err, "failed to create storage: %v", err)

	filename := "test.txt"
	content := []byte("some content")

	err = afero.WriteFile(afs, filename, content, 0777)
	require.NoError(t, err, "failed to create test file: %v", err)

	reader, err := storage.RetrieveFile(filename)
	require.NoError(t, err, "expected to retrieve file without error: %v", err)

	retrievedContent, err := afero.ReadAll(reader)
	require.NoError(t, err, "failed to read retrieved file content: %v", err)

	assert.Equal(
		t,
		content,
		retrievedContent,
		"retrieved file content does not match original content",
	)
}

func TestGetFileUploadsLocation(t *testing.T) {
	afs := afero.NewMemMapFs()
	storage, err := NewDiskFileStorage("uploads", afs)
	require.NoError(t, err, "failed to create storage: %v", err)
	assert.Equal(t, "uploads", storage.GetFileUploadsLocation())

	wd, err := os.Getwd()
	require.NoError(t, err, "failed to get working directory: %v", err)

	path := filepath.Join(wd, "uploads")
	storage, err = NewDiskFileStorage(path, afs)
	require.NoError(t, err, "failed to create storage: %v", err)
	assert.Equal(t, path, storage.GetFileUploadsLocation())
}

func TestShutdown(t *testing.T) {
	afs := afero.NewMemMapFs()
	storage, err := NewDiskFileStorage("uploads", afs)
	require.NoError(t, err)

	err = storage.Shutdown()
	assert.NoError(t, err)
}
