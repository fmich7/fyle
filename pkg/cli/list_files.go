package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/fmich7/fyle/pkg/server"
	"github.com/spf13/cobra"
)

// NewListTreeCmd lists all the files that user stores on the server.
func (c *CliClient) NewListTreeCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "ls [serverPath]",
		Short: "List all your files stored on fyle (as tree).",
		Run: func(cmd *cobra.Command, args []string) {
			path := ""
			if len(args) >= 1 {
				path = args[0]
			}
			if err := c.ListFiles(path); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "error: %v", err)
			}
		},
	}
}

// ListFiles makes request to list user's storage file structure.
func (c *CliClient) ListFiles(path string) error {
	body := new(bytes.Buffer)
	if err := json.NewEncoder(body).Encode(server.ListFilesRequest{Path: path}); err != nil {
		return errors.New("failed to create request body")
	}

	// jwt from keyring
	jwtTokenBytes, err := c.getKeyringValue("jwt_token")
	if err != nil {
		return errors.New("failed to get authorization credentials")
	}
	jwtToken := string(jwtTokenBytes)

	req, err := http.NewRequest("POST", c.ListFilesURL, body)
	if err != nil {
		return errors.New("failed to construct request")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: c.RequestTimeoutTime,
	}

	res, err := client.Do(req)
	if err != nil {
		return errors.New("impossible to send a request")
	}
	defer res.Body.Close()

	msg, err := io.ReadAll(res.Body)
	if err != nil {
		return errors.New("failed to read server response")
	}

	// Is resp good?
	if res.StatusCode != http.StatusOK {
		return errors.New(string(msg))
	}

	_, err = fmt.Fprintf(c.rootCmd.OutOrStdout(), "Your storage tree:\n%s", string(msg))
	return err
}
