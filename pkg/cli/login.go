package cli

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/fmich7/fyle/pkg/crypto"
	"github.com/fmich7/fyle/pkg/server"
	"github.com/spf13/cobra"
)

// NewLogin creates a new login command.
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

// LoginUser sends credentials to the server to login and store JWT token.
func (c *CliClient) LoginUser(username, password string) error {
	if username == "" {
		return errors.New("username argument is empty")
	} else if password == "" {
		return errors.New("password argument is empty")
	}

	data, err := c.makeLoginRequest(username, password)
	if err != nil {
		return err
	}

	loginCredentials := new(server.LoginResponse)
	if err = json.Unmarshal(data, loginCredentials); err != nil {
		return errors.New("failed to unmarshal server response")
	}

	// decode salt
	salt, err := base64.StdEncoding.DecodeString(loginCredentials.Salt)
	if err != nil {
		return errors.New("failed to decode salt")
	}

	// store jwt token
	err = c.setKeyringValue("jwt_token", []byte(loginCredentials.Token))
	if err != nil {
		return err
	}

	// generate symmetric key from password and salt
	encryptionKey := crypto.GeneratePBEKey(password, salt)

	// store encryption key
	err = c.setKeyringValue("encryption_key", encryptionKey)
	if err != nil {
		return err
	}

	return nil
}

// makeLoginRequest makes login request to the server.
func (c *CliClient) makeLoginRequest(username, password string) ([]byte, error) {
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(server.AuthUserRequest{
		Username: username,
		Password: password,
	})
	if err != nil {
		return nil, errors.New("failed to create request body")
	}

	req, err := http.NewRequest("GET", c.LoginURL, body)
	if err != nil {
		return nil, errors.New("failed to construct request")
	}

	client := http.Client{
		Timeout: c.RequestTimeoutTime,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.New("impossible to send a request")
	}
	defer res.Body.Close()

	msg, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, errors.New("failed to read server response")
	}

	// Is resp good?
	if res.StatusCode != http.StatusOK {
		return nil, errors.New(string(msg))
	}

	return msg, nil
}
