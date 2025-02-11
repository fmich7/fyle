package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fyle",
	Short: "fyle is a cli tool for managing your files on the cloud",
}

const (
	uploadURL = "http://localhost:3000/upload"
	user      = "fmich7"
	location  = "folder"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error on executing cli: %s\n", err)
	}
}
