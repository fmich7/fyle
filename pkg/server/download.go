package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/fmich7/fyle/pkg/utils"
)

// HandleFileDownload handles file download request.
func (s *Server) HandleFileDownload(w http.ResponseWriter, r *http.Request) {
	log.Println("Downloading file")

	username := r.Context().Value("username").(string)
	var reqBody DownloadRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "decoding request body", http.StatusBadRequest)
		return
	}

	// filename from req path
	filename := utils.GetFileNameFromPath(reqBody.Path)

	// does file exist on the server?
	path, valid := utils.GetLocationOnServer(
		s.store.GetFileUploadsLocation(),
		username,
		reqBody.Path,
		"", // path already contains a filename
	)

	if !valid {
		http.Error(w, "Provided path is not valid", http.StatusBadRequest)
		return
	}

	// open file stream
	fileReader, err := s.store.RetrieveFile(path)
	if err != nil {
		http.Error(w, "File does not exist", http.StatusBadRequest)
		return
	}
	defer fileReader.Close()

	w.Header().Set("Content-Disposition", fmt.Sprintf(
		"attachment; filename=%s", filename,
	))
	w.Header().Set("Content-Type", "application/octet-stream")

	// send file over the wire
	if _, err := io.Copy(w, fileReader); err != nil {
		http.Error(w, "streaming file", http.StatusInternalServerError)
		return
	}
}
