package server

import (
	"encoding/json"
	"net/http"

	"github.com/fmich7/fyle/pkg/types"
)

// HandleListFiles returns user file tree.
func (s *Server) HandleListFiles(w http.ResponseWriter, r *http.Request) {
	username := r.Context().Value("username").(string)

	var reqBody types.ListFilesRequest
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		http.Error(w, "decoding request body", http.StatusBadRequest)
		return
	}

	// get user file tree
	output, err := s.store.GetUserFileTree(username, reqBody.Path)
	if err != nil {
		http.Error(w, "getting file tree", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}
