package cli

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

// path on server
// filepath - file on disk
var uploadCmd = &cobra.Command{
	Use: "upload",
	Short: "Uploads a file to server\n" +
		"As of now, only single file upload is supported\n" +
		"Usage: fyle upload <path>",
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		filepath := args[0]
		path := "."
		if len(args) > 1 {
			path = args[1]
		}

		UploadFile(filepath, path)
	},
}

// UploadFile uploads a file to the server
func UploadFile(filepath, location string) error {
	file, err := os.Open(filepath)
	if err != nil {
		return fmt.Errorf("Error opening file: %s", err)
	}
	defer file.Close()

	// Create a buffer to store multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	if err := writer.WriteField("user", User); err != nil {
		return err
	}
	if err := writer.WriteField("path", location); err != nil {
		return err
	}

	// Create form
	part, err := writer.CreateFormFile("file", filepath)
	if err != nil {
		return fmt.Errorf("Error creating form: %s", err)
	}
	_, err = io.Copy(part, file)
	if err != nil {
		return fmt.Errorf("Error copying file to form: %s", err)
	}

	// Close writer to finalize form
	writer.Close()

	// Create request and set headers
	resp, err := http.Post(UploadURL, writer.FormDataContentType(), body)
	if err != nil {
		return fmt.Errorf("Error creating request: %s", err)
	}
	defer resp.Body.Close()

	// Check if request was successful
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Error uploading file: %s", resp.Status)
	}

	// Print response
	respBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Server Response:", string(respBody))

	return nil
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
