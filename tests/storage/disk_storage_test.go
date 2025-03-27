package storage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fmich7/fyle/pkg/file"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDiskStorage(t *testing.T) {
	storage, err := storage.NewDiskFileStorage("uploads", afero.NewMemMapFs())
	require.NoError(t, err, "failed to create storage: %v", err)

	wd, err := os.Getwd()
	require.NoError(t, err, "failed to get working directory: %v", err)

	assert.Equal(
		t,
		filepath.Join(wd, "uploads"),
		storage.GetFileUploadsLocation(),
		"unexpected file upload location",
	)
}

func TestStoreFile(t *testing.T) {
	afs := afero.NewMemMapFs()
	storage, err := storage.NewDiskFileStorage("uploads", afs)
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
	storage, err := storage.NewDiskFileStorage("uploads", afs)
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
	storage, err := storage.NewDiskFileStorage("uploads", afs)
	require.NoError(t, err, "failed to create storage: %v", err)

	wd, err := os.Getwd()
	require.NoError(t, err, "failed to get working directory: %v", err)

	expected := filepath.Join(wd, "uploads")
	assert.Equal(t, expected, storage.GetFileUploadsLocation())
}
