package api

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/fmich7/fyle/pkg/utils"
)

// HandleFileDownload handles the file download request
func (s *Server) HandleFileDownload(w http.ResponseWriter, r *http.Request) {
	log.Println("Downloading file")

	// TODO: Auth check

	// Get user from the request
	user := r.FormValue("user")
	if user == "" {
		http.Error(w, "User not provided", http.StatusBadRequest)
		return
	}

	// Get file path from the request
	path := r.FormValue("path")
	if path == "" {
		http.Error(w, "Path not provided", http.StatusBadRequest)
		return
	}

	// Get filename from valid path
	filename := utils.GetFileNameFromPath(path)

	// Check if file exists on a server
	path, valid := utils.GetLocationOnServer(
		s.store.GetFileUploadsLocation(),
		user,
		path,
		"", // path should already contain a filename
	)

	if !valid {
		http.Error(w, "File doesn't exist on a server", http.StatusBadRequest)
	}

	// Open file stream
	fileReader, err := s.store.RetrieveFile(path)
	if err != nil {
		http.Error(w, "Error reading file on a server", http.StatusInternalServerError)
	}
	defer fileReader.Close()

	// Send file over http
	w.Header().Set("Content-Disposition", fmt.Sprintf(
		"attachment; filename=%s", filename,
	))
	w.Header().Set("Content-Type", "application/octet-stream")

	if _, err := io.Copy(w, fileReader); err != nil {
		http.Error(w, "Error streaming file", http.StatusInternalServerError)
	}
}
