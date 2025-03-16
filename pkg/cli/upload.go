package cli

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/fmich7/fyle/pkg/encryption"
	"github.com/fmich7/fyle/pkg/types"
	"github.com/spf13/cobra"
)

// NewUploadCmd creates a new upload command
func (c *CliClient) NewUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "upload [localPath] [serverPath]",
		Short: "Uploads a file to server from given location",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			localPath := args[0]
			serverPath := "."
			if len(args) > 1 {
				serverPath = args[1]
			}

			if err := c.UploadFile(localPath, serverPath); err != nil {
				fmt.Fprintf(cmd.ErrOrStderr(), "error uploading file: %v", err)
			}
		},
	}
}

// UploadFile uploads a file to the server
func (c *CliClient) UploadFile(localPath, serverPath string) error {
	// Get stored JWT token and set it as Authorization header
	jwtTokenBytes, err := c.getKeyringValue("jwt_token")
	if err != nil {
		return errors.New("failed to get authorization credentials")
	}
	jwtToken := string(jwtTokenBytes)

	// Get encryption key from keyring
	encryptionKey, err := c.getKeyringValue("encryption_key")
	_ = encryptionKey
	if err != nil {
		return errors.New("failed to get encryption_key")
	}

	// Prepare multipart form data for the request
	// It doesn't load the file into memory
	form, err := c.PrepareMultipartForm(localPath, serverPath, encryptionKey)
	if err != nil {
		return err
	}

	// Create request and set headers
	req, err := http.NewRequest("POST", c.UploadURL, form.FormData)
	req.Header.Set("Content-Type", form.FormDataContentType)
	if err != nil {
		return fmt.Errorf("creating request: %v", err)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))

	client := http.Client{
		Timeout: c.RequestTimeoutTime,
	}

	// Send request
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("impossible to send a request: %v", err)
	}

	defer res.Body.Close()

	// Check if request was successful
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("uploading file: %v", res.Status)
	}

	// Print response
	respBody, _ := io.ReadAll(res.Body)
	fmt.Println("Server Response:", string(respBody))

	return nil
}

// PrepareMultipartForm prepares a multipart form for the request
// It doesn't load the file into memory
func (c *CliClient) PrepareMultipartForm(
	localPath, serverPath string, encryptionKey []byte,
) (*types.MultiPartForm, error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	go func() {
		defer w.Close()
		defer m.Close()

		// write the path field
		if err := m.WriteField("path", serverPath); err != nil {
			w.CloseWithError(fmt.Errorf("writing path field: %w", err))
			return
		}

		// create a form file
		formFile, err := m.CreateFormFile("file", localPath)
		if err != nil {
			w.CloseWithError(fmt.Errorf("creating form file: %w", err))
			return
		}

		// open the local file
		file, err := c.fs.Open(localPath)
		if err != nil {
			w.CloseWithError(fmt.Errorf("opening file: %w", err))
			return
		}
		defer file.Close()

		encryptedFileReader, err := encryption.EncryptData(file, encryptionKey)
		if err != nil {
			w.CloseWithError(fmt.Errorf("encrypting file: %w", err))
			return
		}

		// copy the file to the form part
		if _, err := io.Copy(formFile, encryptedFileReader); err != nil {
			w.CloseWithError(fmt.Errorf("copying file to form: %w", err))
			return
		}
	}()

	return &types.MultiPartForm{
		FormData:            r,
		FormDataContentType: m.FormDataContentType(),
	}, nil
}
