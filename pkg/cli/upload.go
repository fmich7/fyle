package cli

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/fmich7/fyle/pkg/types"
	"github.com/spf13/afero"
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
	file, err := c.fs.Open(localPath)
	if err != nil {
		return fmt.Errorf("opening file: %v", err)
	}
	defer file.Close()

	// Create a buffer to store multipart form data
	form, err := PrepareMultipartForm(file, localPath, serverPath)
	if err != nil {
		return err
	}

	// Create request and set headers
	req, err := http.NewRequest("POST", c.UploadURL, form.FormData)
	req.Header.Set("Content-Type", form.FormDataContentType)
	if err != nil {
		return fmt.Errorf("creating request: %v", err)
	}

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
	// Check if request was successful
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("uploading file: %v", res.Status)
	}

	// Print response
	respBody, _ := io.ReadAll(res.Body)
	fmt.Println("Server Response:", string(respBody))

	return nil
}

func PrepareMultipartForm(file afero.File, localPath, serverPath string) (*types.MultiPartForm, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	if err := writer.WriteField("path", serverPath); err != nil {
		return nil, err
	}

	// Create form
	part, err := writer.CreateFormFile("file", localPath)
	if err != nil {
		return nil, fmt.Errorf("creating form: %v", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, fmt.Errorf("copying file to form: %v", err)
	}

	return &types.MultiPartForm{
		FormData:            body,
		FormDataContentType: writer.FormDataContentType(),
	}, nil
}
