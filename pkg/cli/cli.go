package cli

import (
	"fmt"
	"os"
	"time"

	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// CliClient struct client type for the cli.
type CliClient struct {
	rootCmd            *cobra.Command
	fs                 afero.Fs
	UploadURL          string
	DownloadURL        string
	LoginURL           string
	ListFilesURL       string
	SignupURL          string
	KeyRingName        string
	RequestTimeoutTime time.Duration
}

// NewCliClient creates a new CliClient object with specified file system.
func NewCliClient(fs afero.Fs) *CliClient {
	client := &CliClient{
		rootCmd: &cobra.Command{
			Use:   "fyle",
			Short: "fyle is a cli tool for managing your files on the cloud",
		},
		fs: fs,

		UploadURL:          "http://localhost:3000/file",
		DownloadURL:        "http://localhost:3000/getfile",
		LoginURL:           "http://localhost:3000/login",
		SignupURL:          "http://localhost:3000/signup",
		ListFilesURL:       "http://localhost:3000/ls",
		KeyRingName:        "fyle",
		RequestTimeoutTime: 10 * time.Second,
	}

	client.attachCommands()

	return client
}

// Execute runs the root command.
func (c *CliClient) Execute() {
	if err := c.rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "error executing cli: %v", err)
	}
}

// Attaches commands to the client.
func (c *CliClient) attachCommands() {
	c.rootCmd.AddCommand(c.NewUploadCmd())
	c.rootCmd.AddCommand(c.NewDownloadCmd())
	c.rootCmd.AddCommand(c.NewLoginCmd())
	c.rootCmd.AddCommand(c.NewSignUPCmd())
	c.rootCmd.AddCommand(c.NewListTreeCmd())
}
