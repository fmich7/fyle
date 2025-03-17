package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fmich7/fyle/pkg/encryption"
	"github.com/fmich7/fyle/pkg/types"
	"github.com/fmich7/fyle/pkg/utils"
	"github.com/spf13/cobra"
)

// NewDownloadCmd creates a new download command
func (c *CliClient) NewDownloadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "download [serverPath] [localPath]",
		Short: "Downloads a file from server and stores it in given location",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serverPath := args[0]
			localPath := "."
			if len(args) > 1 {
				localPath = args[1]
			}

			localPath, err := utils.GetBaseDir(localPath)
			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), "error: couldn't process given path")
				return
			}

			if err = c.DownloadFile(serverPath, localPath); err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), fmt.Errorf("error: %v", err))
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "File saved at %s", localPath)
			}
		},
	}
}

// DownloadFile handles download request from cli
func (c *CliClient) DownloadFile(serverPath, localPath string) error {
	data := types.DownloadRequest{
		Path: serverPath,
	}

	// Marshall request data
	marshalled, err := json.Marshal(data)
	if err != nil {
		return errors.New("marshalling data")
	}

	// Send request
	req, err := http.NewRequest("POST", c.DownloadURL, bytes.NewBuffer(marshalled))
	if err != nil {
		return errors.New("couldn't construct a request")
	}
	req.Header.Set("Content-Type", "application/json")

	// Load jwt token from keyring and set it as Authorization header
	jwtTokenBytes, err := c.getKeyringValue("jwt_token")
	if err != nil {
		return errors.New("failed to get authorization credentials")
	}
	jwtToken := string(jwtTokenBytes)
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))

	client := http.Client{
		Timeout: c.RequestTimeoutTime,
	}

	// Send request
	res, err := client.Do(req)
	if err != nil {
		return errors.New("impossible to send a request")
	}
	defer res.Body.Close()

	// Check disposition header and get filename
	dispositionHeader := res.Header.Get("Content-Disposition")
	filename, err := utils.GetFileNameFromContentDisposition(dispositionHeader)
	if err != nil {
		return errors.New("failed to get filename")
	}

	// Get encryption key from keyring
	encryptionKey, err := c.getKeyringValue("encryption_key")
	if err != nil {
		return errors.New("failed to get encryption_key")
	}

	// Decryption stream from response body
	decryptionFileStream := encryption.DecryptData(res.Body, encryptionKey)

	// Save file on disk
	err = utils.SaveFileOnDisk(c.fs, localPath, filename, decryptionFileStream)
	if err != nil {
		return fmt.Errorf("failed to save file on disk %w", err)
	}

	return nil
}
