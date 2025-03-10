package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "fyle",
	Short: "fyle is a cli tool for managing your files on the cloud",
}

// TODO: CONFIG!!!!!
var (
	UploadURL          = "http://localhost:3000/upload"
	DownloadURL        = "http://localhost:3000/download"
	User               = "fmich7"
	RequestTimeoutTime = 10 * time.Second
)

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error on executing cli: %v", err)
	}
}
