package file

import (
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetDirTree(t *testing.T) {
	afs := afero.NewMemMapFs()
	err := afs.MkdirAll("/test/folder", os.ModePerm)
	require.NoError(t, err)
	err = afero.WriteFile(afs, "/test/file", []byte(""), os.ModePerm)
	require.NoError(t, err)

	res, err := GetDirTree(afs, "/")
	require.NoError(t, err)

	expected := "|── ./\n|    |── test/\n|    |    |── file\n|    |    |── folder/\n"
	assert.Equal(t, expected, res)
}
