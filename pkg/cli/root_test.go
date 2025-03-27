package cli

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/assert"
)

func TestExecute_NoError(t *testing.T) {
	output := &bytes.Buffer{}
	rootCmd := &cobra.Command{
		Use:   "fyle",
		Short: "fyle is a cli tool for managing your files on the cloud",
	}
	rootCmd.SetOut(output)

	err := rootCmd.Execute()
	assert.NoError(t, err, "Executing root command should not return an error")
}
