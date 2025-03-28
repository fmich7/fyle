package server

import (
	"testing"
	"time"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/fmich7/fyle/pkg/storage"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestServerStart(t *testing.T) {
	cfg := config.NewTestingConfig()
	afs := afero.NewMemMapFs()
	store, err := storage.NewTestingStorage(afs)
	require.NoError(t, err)

	srv := NewServer(cfg, store)
	assert.False(t, srv.IsRunning(), "IsRunning should be false before start")

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

	assert.True(t, srv.IsRunning(), "IsRunning should be true after start")

	port, err := srv.GetPort()
	require.NoError(t, err, "Failed to get server port")
	require.Greater(t, port, 0, "Invalid port assigned")
	require.NoError(t, srv.Shutdown(), "Failed to shutdown server")
}
