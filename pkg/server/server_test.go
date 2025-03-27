package server

import (
	"testing"
	"time"

	"github.com/fmich7/fyle/pkg/config"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	s := NewServer(config.NewTestingConfig(), nil)
	errCh := make(chan error, 1)
	defer close(errCh)

	go func() {
		errCh <- s.Start()
	}()

	select {
	case err := <-errCh:
		assert.NoError(t, err, "Expected no error, got: %v", err)
	case <-time.After(10 * time.Millisecond):
		break
	}
}
