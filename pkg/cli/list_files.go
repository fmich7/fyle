package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/fmich7/fyle/pkg/types"
	"github.com/spf13/cobra"
)

// NewListTreeCmd lists all files that user has stored
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

// ListFiles prints user's files stored on server
func (c *CliClient) ListFiles(path string) error {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(types.ListFilesRequest{
		Path: path,
	})
	if err != nil {
		return errors.New("failed to create request body")
	}

	req, err := http.NewRequest("POST", c.ListFilesURL, body)
	if err != nil {
		return errors.New("failed to construct request")
	}
	req.Header.Set("Content-Type", "application/json")

	jwtToken, err := c.getJWTToken()
	if err != nil {
		return errors.New("failed to get authorization credentials")
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))

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

	fmt.Printf("Your storage tree:\n%s", string(msg))

	return nil
}
