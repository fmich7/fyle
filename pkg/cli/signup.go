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

// NewSignUPCmd creates a new sign up command.
func (c *CliClient) NewSignUPCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "signup [username] [password]",
		Short: "Create a new account on fyle platform",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			password := args[1]
			if err := c.SignUPUser(username, password); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "error: %v", err)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "Created new account successfully!")
			}
		},
	}
}

// SignUPUser sends singup request to the server.
func (c *CliClient) SignUPUser(username, password string) error {
	if username == "" {
		return errors.New("username argument is empty")
	} else if password == "" {
		return errors.New("password argument is empty")
	}

	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(types.AuthUserRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return errors.New("failed to create request body")
	}

	req, err := http.NewRequest("POST", c.SignupURL, body)
	if err != nil {
		return errors.New("failed to construct request")
	}

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
	if res.StatusCode != http.StatusCreated {
		return errors.New(string(msg))
	}

	return nil
}
