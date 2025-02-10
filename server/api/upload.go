package api

import (
	"log"
	"net/http"

	"github.com/fmich7/fyle/types"
)

// HandleFileUpload handles the file upload request
func (s *Server) HandleFileUpload(w http.ResponseWriter, r *http.Request) {
	log.Println("Uploading file")
	r.ParseMultipartForm(10 << 20) // 10 MB max size

	// Get user and location from the request
	// TODO: AUTH CHECK
	if !r.Form.Has("user") {
		http.Error(w, "Missing user", http.StatusBadRequest)
		return
	}
	user := r.FormValue("user")
	location := r.FormValue("location")

	// Retrieve the file from the request
	fileData, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusBadRequest)
		return
	}
	defer fileData.Close()

	// Create a new file object and upload it
	file := types.NewFile(header, fileData, user, location)
	if err := s.store.UploadFile(file); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Return a 201 Created status
	w.WriteHeader(http.StatusCreated)
	log.Println("File uploaded successfully: ", file.Filename)
}
