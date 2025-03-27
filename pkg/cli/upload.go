package cli

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/fmich7/fyle/pkg/file"
	"github.com/fmich7/fyle/pkg/utils"
	"github.com/spf13/cobra"
)

// NewUploadCmd creates a new upload command.
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

// UploadFile sends file to the server.
func (c *CliClient) UploadFile(localPath, serverPath string) error {
	jwtTokenBytes, err := c.getKeyringValue("jwt_token")
	if err != nil {
		return errors.New("failed to get authorization credentials")
	}
	jwtToken := string(jwtTokenBytes)

	encryptionKey, err := c.getKeyringValue("encryption_key")
	if err != nil {
		return errors.New("failed to get encryption_key")
	}

	// Prepare multipart form data for the request
	// It doesn't load the file into memory
	form, err := c.PrepareMultipartForm(localPath, serverPath, encryptionKey)
	if err != nil {
		return err
	}

	// create request
	req, err := http.NewRequest("POST", c.UploadURL, form.FormData)
	if err != nil {
		return fmt.Errorf("creating request: %v", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwtToken))
	req.Header.Set("Content-Type", form.FormDataContentType)

	client := http.Client{
		Timeout: c.RequestTimeoutTime,
	}

	// send erquest
	res, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("impossible to send a request: %v", err)
	}

	defer res.Body.Close()

	msg, _ := io.ReadAll(res.Body)
	// goood rqeuest>?
	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("uploading file: %s", string(msg))
	}

	fmt.Println("Server Response:", string(msg))

	return nil
}

// PrepareMultipartForm writes filepath, and file data in chunks to a multipart request.
func (c *CliClient) PrepareMultipartForm(
	localPath, serverPath string, encryptionKey []byte,
) (*file.MultiPartForm, error) {
	r, w := io.Pipe()
	m := multipart.NewWriter(w)

	go func() {
		defer w.Close()

		// write filepath on the server
		if err := m.WriteField("path", serverPath); err != nil {
			w.CloseWithError(fmt.Errorf("writing path field: %v", err))
			return
		}

		// create form file
		formFile, err := m.CreateFormFile("file", localPath)
		if err != nil {
			w.CloseWithError(fmt.Errorf("creating form file: %v", err))
			return
		}

		file, err := c.fs.Open(localPath)
		if err != nil {
			w.CloseWithError(fmt.Errorf("opening file: %v", err))
			return
		}
		defer file.Close()

		// encrypt file in chunks (stream)
		encryptionFileStream := utils.EncryptData(file, encryptionKey)

		if _, err := io.Copy(formFile, encryptionFileStream); err != nil {
			w.CloseWithError(fmt.Errorf("error copying encrypted data: %v", err))
			return
		}

		if err := m.Close(); err != nil {
			w.CloseWithError(fmt.Errorf("error closing multipart writer: %v", err))
			return
		}
	}()

	return &file.MultiPartForm{
		FormData:            r,
		FormDataContentType: m.FormDataContentType(),
	}, nil
}
