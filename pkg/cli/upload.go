package cli

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/fmich7/fyle/pkg/utils"
	"github.com/spf13/cobra"
)

// UploadFile uploads a file to the server
func UploadFile(path string) error {
	absPath, err := utils.GetAbsolutePath(path)
	if err != nil {
		return err
	}

	// Open file
	file, err := os.Open(absPath)
	if err != nil {
		return fmt.Errorf("Error opening file: %s", err)
	}
	defer file.Close()

	// Create a buffer to store multipart form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	defer writer.Close()

	if err := writer.WriteField("user", user); err != nil {
		return err
	}
	if err := writer.WriteField("location", location); err != nil {
		return err
	}

	// Create form
	part, err := writer.CreateFormFile("file", absPath)
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
	resp, err := http.Post(uploadURL, writer.FormDataContentType(), body)
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

// TODO: Add multiple file upload support
var uploadCmd = &cobra.Command{
	Use: "upload",
	Short: "Uploads a file to server\n" +
		"As of now, only single file upload is supported\n" +
		"Usage: fyle upload <file-path>",
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		UploadFile(args[0])
	},
}

func init() {
	rootCmd.AddCommand(uploadCmd)
}
