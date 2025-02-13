package api

import (
	"fmt"
	"log"
	"net/http"

	"github.com/fmich7/fyle/pkg/types"
	"github.com/fmich7/fyle/pkg/utils"
)

// HandleFileUpload handles the file upload request
func (s *Server) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("Uploading file")
	r.ParseMultipartForm(10 << 20) // 10 MB max size

	// Get user and location from the request
	// TODO: AUTH CHECK

	user := r.FormValue("user")
	if user == "" {
		http.Error(w, "User not provided", http.StatusBadRequest)
		fmt.Println("User not provided")
		return
	}

	// Retrieve the file from the request
	fileData, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		fmt.Println("Error retrieving file")
		return
	}
	defer fileData.Close()

	userInputPath := r.FormValue("location")
	safePath, valid := utils.LocationOnServer(
		s.store.GetFileUploadsLocation(),
		user,
		userInputPath,
		header.Filename,
	)

	if !valid {
		http.Error(w, "Invalid location", http.StatusBadRequest)
		fmt.Println(safePath, userInputPath, user, s.store.GetFileUploadsLocation(), header.Filename)
		return
	}

	// Create a new file object and upload it
	file := types.NewFile(header, fileData, user, safePath)
	if err := s.store.UploadFile(file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		fmt.Println("Error uploading file:", err)
		return
	}

	// Return a 201 Created status
	w.WriteHeader(http.StatusCreated)
	log.Println("File uploaded successfully:", file.Owner, file.Filename)
}
