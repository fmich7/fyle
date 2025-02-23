package cli

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/fmich7/fyle/pkg/types"
	"github.com/fmich7/fyle/pkg/utils"
	"github.com/spf13/cobra"
)

// DownloadFunc is a function that downloads a file from a server
// It's purpose is to easily replace this in tests
type DownloadFunc func(string, string) error

// NewDownloadCmd creates a new download command
// serverPath is the path to the file on the server
// destination is the path to the directory where the file will be saved
func NewDownloadCmd(download DownloadFunc) *cobra.Command {
	return &cobra.Command{
		Use: "download",
		Short: "Downloads a file from server\n" +
			"As of now, only single file download is supported\n" +
			"Usage: fyle download <path> <destination>",
		Args: cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			serverPath := args[0]
			destination := "."
			if len(args) > 1 {
				destination = args[1]
			}

			destination, err := utils.GetBaseDir(destination)
			if err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), "error: couldn't process given path")
				return
			}

			if err = download(serverPath, destination); err != nil {
				fmt.Fprintln(cmd.ErrOrStderr(), err)
			} else {
				fmt.Fprintf(cmd.OutOrStdout(), "File saved at %s\n", destination)
			}
		},
	}
}

// DownloadFile handles download request from cli
func DownloadFile(serverPath, destination string) error {
	data := types.DownloadRequest{
		Path: serverPath,
		// TODO: AUTH!!!!!!
		User: "fmich7",
	}

	// Marshall request data
	marshalled, err := json.Marshal(data)
	if err != nil {
		return errors.New("error: marshalling data")
	}

	// Send request
	req, err := http.NewRequest("POST", DownloadURL, bytes.NewBuffer(marshalled))
	if err != nil {
		return errors.New("error: couldn't construct a request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: RequestTimeoutTime,
	}

	// Send request
	res, err := client.Do(req)
	if err != nil {
		return errors.New("error: impossible to send a request")
	}
	defer res.Body.Close()

	// Check disposition header and get filename
	dispositionHeader := res.Header.Get("Content-Disposition")
	filename, err := utils.GetFileNameFromContentDisposition(dispositionHeader)
	if err != nil {
		return errors.New("error: bad request headers")
	}

	// Save file on disk
	err = utils.SaveFileOnDisk(destination, filename, res.Body)
	if err != nil {
		return errors.New("error: couldn't save file on disk")
	}

	return nil
}

func init() {
	rootCmd.AddCommand(NewDownloadCmd(DownloadFile))
}
