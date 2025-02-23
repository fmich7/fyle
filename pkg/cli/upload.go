package cli

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/fmich7/fyle/pkg/types"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

// NewUploadCmd creates a new upload command
func (c *CliClient) NewUploadCmd() *cobra.Command {
	return &cobra.Command{
		Use: "upload",
		Short: "Uploads a file to server\n" +
			"As of now, only single file upload is supported\n" +
			"Usage: fyle upload <path>",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			localPath := args[0]
			serverPath := "."
			if len(args) > 1 {
				serverPath = args[1]
			}

			if err := c.UploadFile(localPath, serverPath); err != nil {
				fmt.Fprintf(os.Stderr, "error uploading file: %v\n", err)
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
	form, err := PrepareMultipartForm(file, localPath, serverPath, c.User)
	if err != nil {
		return err
	}

	// Create request and set headers
	resp, err := http.Post(c.UploadURL, form.FormDataContentType, form.FormData)
	if err != nil {
		return fmt.Errorf("creating request: %v", err)
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("uploading file: %v", resp.Status)
	}

	// Print response
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Server Response:", string(respBody))

	return nil
}

func PrepareMultipartForm(file afero.File, localPath, serverPath, user string) (*types.MultiPartForm, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	if err := writer.WriteField("user", user); err != nil {
		return nil, err
	}
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
