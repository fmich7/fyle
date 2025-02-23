package server

import (
	"log"
	"net/http"

	"github.com/fmich7/fyle/pkg/types"
	"github.com/fmich7/fyle/pkg/utils"
)

// HandleFileUpload handles the file upload request
func (s *Server) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("Uploading file")
	r.ParseMultipartForm(10 << 20) // 10 MB max size

	// Get user from the request
	user := r.FormValue("user")
	if user == "" {
		http.Error(w, "User not provided", http.StatusBadRequest)
		return
	}

	// Retrieve the file from the request
	fileData, fileMetadata, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "error retrieving file", http.StatusBadRequest)
		return
	}
	defer fileData.Close()

	// Check if requested path is valid
	userInputPath := r.FormValue("path")
	safePath, valid := utils.GetLocationOnServer(
		s.store.GetFileUploadsLocation(),
		user,
		userInputPath,
		fileMetadata.Filename,
	)

	if !valid {
		http.Error(w, "Invalid location", http.StatusBadRequest)
		return
	}

	// Create new file object and upload it to the storage
	file := types.NewFile(fileMetadata, fileData, user, safePath)
	if err := s.store.StoreFile(file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 201 Created status
	w.WriteHeader(http.StatusCreated)
	log.Println("File uploaded successfully:", file.Owner, file.Filename)
}
