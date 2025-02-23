package cli

import (
	"bytes"
	"io"
	"os"
	"testing"

	"github.com/fmich7/fyle/pkg/cli"
	"github.com/spf13/afero"
)

func TestExecute(t *testing.T) {
	// Create a pipe to capture the output
	r, w, _ := os.Pipe()
	oldStderr := os.Stderr
	os.Stderr = w
	defer func() {
		os.Stderr = oldStderr
		w.Close()
	}()

	// Test function
	cli := cli.NewCliClient(afero.NewMemMapFs())
	cli.Execute()

	// Close the write end of the pipe and read the output
	if err := w.Close(); err != nil {
		t.Errorf("failed to close pipe: %v\n", err)
	}

	var buf bytes.Buffer
	io.Copy(&buf, r)

	// Check if there was an error
	if buf.Len() > 0 {
		t.Errorf("Expected no error, got: %s", buf.String())
	}
}
