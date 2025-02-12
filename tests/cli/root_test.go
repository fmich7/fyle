package cli_test

import (
	"bytes"
	"os"
	"testing"

	"github.com/fmich7/fyle/pkg/cli"
)

func TestExecute(t *testing.T) {
	// Create a pipe to retransfer the output
	r, w, _ := os.Pipe()
	oldStderr := os.Stderr
	os.Stderr = w
	defer func() {
		os.Stderr = oldStderr
		w.Close()
	}()

	// Test function
	cli.Execute()

	// Close the write end of the pipe and read the output
	w.Close()

	var buf bytes.Buffer
	buf.ReadFrom(r)

	// Check if there was an error
	if buf.Len() > 0 {
		t.Errorf("Expected no error, got: %s", buf.String())
	}
}
