package server

import (
	"encoding/json"
	"net/http"
)

// HandleListFiles returns user file tree.
func (s *Server) HandleListFiles(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value(CtxUsernameKey{}).(string)

	var reqBody ListFilesRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "decoding request body", http.StatusBadRequest)
		return
	}

	// make user absolute path to the target directory
	fullPath := JoinPathParts(s.store.GetFileUploadsLocation(), username, reqBody.Path)
	valid := ValidatePath(s.store.GetFileUploadsLocation(), fullPath)

	if !valid {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	// get user file tree
	output, err := s.store.GetUserFileTree(fullPath)
	if err != nil {
		http.Error(w, "getting file tree", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}
