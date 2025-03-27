package server

import (
	"testing"
	"time"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/require"
)

func TestServerStart(t *testing.T) {
	cfg := config.NewTestingConfig()
	afs := afero.NewMemMapFs()
	store, err := storage.NewTestingStorage(afs)
	require.NoError(t, err)

	srv := NewServer(cfg, store)

	errCh := make(chan error, 1)
	go func() {
		errCh <- srv.Start()
	}()

	select {
	case <-time.After(200 * time.Millisecond):
	case err := <-errCh:
		t.Fatalf("Failed to start server: %v", err)
	}

	port, err := srv.GetPort()
	require.NoError(t, err, "Failed to get server port")
	require.Greater(t, port, 0, "Invalid port assigned")
	require.NoError(t, srv.Shutdown(), "Failed to shutdown server")
}
