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

// NewLogin creates a new login command
func (c *CliClient) NewLoginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "login [username] [password]",
		Short: "Login into your account",
		Args:  cobra.MinimumNArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			username := args[0]
			password := args[1]
			if err := c.LoginUser(username, password); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "error: %v", err)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "Logged in successfully!")
			}
		},
	}
}

// LoginUser sends given credentials to server in order to receive and store JWT
func (c *CliClient) LoginUser(username, password string) error {
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

	req, err := http.NewRequest("GET", c.LoginURL, body)
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
	if res.StatusCode != http.StatusOK {
		return errors.New(string(msg))
	}

	loginCredentials := new(types.LoginResponse)
	err = json.Unmarshal(msg, loginCredentials)
	if err != nil {
		return errors.New("failed to unmarshal server response")
	}

	fmt.Printf("%+v\n", loginCredentials)
	err = c.setJWTToken(loginCredentials.Token)
	if err != nil {
		return err
	}

	return nil
}
