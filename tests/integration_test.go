package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/fmich7/fyle/pkg/cli"
	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/server"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// This tests whole application end-to-end
func TestApplicationIntegrationTest(t *testing.T) {
	// setup server
	cfg := config.NewTestingConfig()
	afs := afero.NewMemMapFs()
	store, err := storage.NewTestingStorage(afs)
	require.NoError(t, err, "Failed to create testing storage")

	srv := server.NewServer(cfg, store)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	// wait for the server to start
	select {
	case <-time.After(100 * time.Millisecond):
		break
	case err := <-errCh:
		t.Fatalf("Failed to start server: %v", err)
	}

	// setup cli
	cli := cli.NewCliClient(afs)
	port, err := srv.GetPort()
	require.NoError(t, err)
	setupCliConfig(cli, port)

	// create an account
	// login to account
	// create file
	// send file
	// retrieve file
	// check if file content matches ;)

	err = srv.Shutdown()
	assert.NoError(t, err, "Failed to graceful shutdown the server")
}

func setupCliConfig(cli *cli.CliClient, port int) {
	cli.UploadURL = fmt.Sprintf("http://localhost:%d/file", port)
	cli.DownloadURL = fmt.Sprintf("http://localhost:%d/getfile", port)
	cli.LoginURL = fmt.Sprintf("http://localhost:%d/login", port)
	cli.SignupURL = fmt.Sprintf("http://localhost:%d/signup", port)
	cli.ListFilesURL = fmt.Sprintf("http://localhost:%d/ls", port)
}
