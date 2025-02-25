package cli

import (
	"fmt"
	"os"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// CliClient is the client type for the cli
type CliClient struct {
	rootCmd     *cobra.Command
	fs          afero.Fs
	UploadURL   string
	DownloadURL string
	User        string
}

// NewCliClient creates a new CliClient object
func NewCliClient(fs afero.Fs) *CliClient {
	client := &CliClient{
		rootCmd: &cobra.Command{
			Use:   "fyle",
			Short: "fyle is a cli tool for managing your files on the cloud",
		},
		fs: fs,

		// TODO: CONFIG!!!!!
		UploadURL:   "http://localhost:3000/file",
		DownloadURL: "http://localhost:3000/getfile",
		User:        "fmich7",
	}

	// Attaches commands to the client
	client.attachCommands()

	return client
}

// Execute runs the root command
func (c *CliClient) Execute() {
	if err := c.rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error executing cli: %v\n", err)
	}
}

func (c *CliClient) attachCommands() {
	c.rootCmd.AddCommand(c.NewUploadCmd())
	c.rootCmd.AddCommand(c.NewDownloadCmd())
}
