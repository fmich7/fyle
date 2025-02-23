package cli

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/fmich7/fyle/pkg/types"
	"github.com/spf13/cobra"
)

var uploadCmd = &cobra.Command{
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

		UploadFile(localPath, serverPath)
	},
}

// UploadFile uploads a file to the server
func UploadFile(localPath, serverPath string) error {
	file, err := os.Open(localPath)
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
	resp, err := http.Post(UploadURL, form.FormDataContentType, form.FormData)
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

func init() {
	rootCmd.AddCommand(uploadCmd)
}

func PrepareMultipartForm(file *os.File, localPath, serverPath string) (*types.MultiPartForm, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	if err := writer.WriteField("user", User); err != nil {
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
