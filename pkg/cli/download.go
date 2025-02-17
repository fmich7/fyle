package cli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/fmich7/fyle/pkg/types"
	"github.com/fmich7/fyle/pkg/utils"
	"github.com/spf13/cobra"
)

// Path is location on a server
// Destination is save location
var downloadCmd = &cobra.Command{
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
			fmt.Fprintln(os.Stderr, "error: couldn't process given path")
			return
		}

		DownloadFile(serverPath, destination)
	},
}

func DownloadFile(serverPath, destination string) {
	data := types.DownloadRequest{
		Path: serverPath,
	}

	marshalled, err := json.Marshal(data)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: marshalling data")
		return
	}

	req, err := http.NewRequest("POST", DownloadURL, bytes.NewBuffer(marshalled))
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: couldn't construct a request")
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{
		Timeout: RequestTimeoutTime,
	}

	res, err := client.Do(req)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: impossible to send a request")
		return
	}
	defer res.Body.Close()

	dispositionHeader := res.Header.Get("Content-Disposition")
	filename, err := utils.GetFileNameFromContentDisposition(dispositionHeader)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: bad request headers")
		return
	}

	err = utils.SaveFileOnDisk(destination, filename, res.Body)
	fmt.Println(destination, filename)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error: couldn't save file on disk", err)
		return
	}
	fmt.Printf("Saved file %s in %s", filename, destination)
}

func init() {
	rootCmd.AddCommand(downloadCmd)
}
