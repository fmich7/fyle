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
	username := r.Context().Value("username").(string)

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
		username,
		userInputPath,
		fileMetadata.Filename,
	)

	if !valid {
		http.Error(w, "Invalid location", http.StatusBadRequest)
		return
	}

	// Create new file object and upload it to the storage
	file := types.NewFile(fileMetadata, fileData, username, safePath)
	if err := s.store.StoreFile(file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return 201 Created status
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("File uploaded successfully"))
	log.Println("File uploaded successfully:", file.Owner, file.Filename)
}
