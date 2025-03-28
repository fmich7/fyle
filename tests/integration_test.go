package tests

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/fmich7/fyle/pkg/cli"
	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/server"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

// This integration tests whole codebase end-to-end.
func TestApplicationIntegrationTest(t *testing.T) {
	// setup server
	cfg := config.NewTestingConfig()
	afs := afero.NewMemMapFs()
	store, err := storage.NewTestingStorage(afs)
	require.NoError(t, err, "Failed to create testing storage")

	// check if uploads folder is created
	exists, err := afero.DirExists(afs, "uploads")
	require.NoError(t, err, "Failed to create testing storage")
	require.True(t, exists)

	srv := server.NewServer(cfg, store)
	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	// wait for the server to start
	timeout := time.After(3 * time.Second)
	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	select {
	case <-ticker.C:
		if srv.IsRunning() {
			break
		}
	case <-timeout:
		t.Fatalf("Server failed to start within the timeout")
	case err := <-errCh:
		t.Fatalf("Failed to start server: %v", err)
	}

	// cleanup even if tests fail
	t.Cleanup(func() { _ = srv.Shutdown() })

	// setup cli
	cli := cli.NewCliClient(afs)
	port, err := srv.GetPort()
	require.NoError(t, err)
	setupCliConfig(cli, port)

	username := "testusername"
	password := "dontstoreplainpasswords"

	// create an account
	err = cli.SignUPUser(username, password)
	require.NoError(t, err, "Failed to create an user from CliClient")

	// login to account
	err = cli.LoginUser(username, password)
	require.NoError(t, err, "Failed to login from CliClient")

	// create file
	filename := "test.data"
	data := []byte("sensitive data!!! 🤯")
	err = afero.WriteFile(afs, filename, data, os.ModePerm)
	require.NoError(t, err, "Failed to create file used for upload")

	// upload file to the server
	err = cli.UploadFile("test.data", "./documents/")
	require.NoError(t, err, "Failed to upload file to the server")

	// check if file exists on the server
	path := fmt.Sprintf("uploads/%s/documents/%s", username, filename)
	exists, err = afero.Exists(afs, path)
	require.True(t, exists, "Uploaded file doesn't exists")
	require.NoError(t, err)

	// check if contents are different (encryption)
	uploadedFileData, err := afero.ReadFile(afs, path)
	require.NoError(t, err, "Failed to read file content")
	require.NotEqual(t, data, uploadedFileData, "File should be encrypted")

	// download file from the server file
	err = afs.Mkdir("downloads", os.ModePerm)
	require.NoError(t, err, "Failed to create downloads dir")
	err = cli.DownloadFile("./documents/test.data", "downloads")
	require.NoError(t, err, "Failed to download file from server")

	// check if file was saved after download
	dwnFilePath := "downloads/test.data"
	exists, err = afero.Exists(afs, dwnFilePath)
	require.NoError(t, err)
	require.True(t, exists, "Downloaded file wasn't saved")

	// doest file content match?
	dwnFileData, err := afero.ReadFile(afs, dwnFilePath)
	require.NoError(t, err, "Failed to read file content")
	require.Equal(t, data, dwnFileData, "File contents doesn't match")
}

func setupCliConfig(cli *cli.CliClient, port int) {
	cli.UploadURL = fmt.Sprintf("http://localhost:%d/file", port)
	cli.DownloadURL = fmt.Sprintf("http://localhost:%d/getfile", port)
	cli.LoginURL = fmt.Sprintf("http://localhost:%d/login", port)
	cli.SignupURL = fmt.Sprintf("http://localhost:%d/signup", port)
	cli.ListFilesURL = fmt.Sprintf("http://localhost:%d/ls", port)
}
