package server

import (
	"errors"
	"io"
	"strings"
)

// LoginResponse contains data that is returned to client after login.
type LoginResponse struct {
	Token string `json:"token"`
	Salt  string `json:"salt"`
}

// DownloadRequest is contains data that server expects from client.
type DownloadRequest struct {
	Path string `json:"path"`
}

// AuthUserRequest contains data that is used to create/login to acc.
type AuthUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ListFilesRequest contains data that is used to list files in user's dir.
type ListFilesRequest struct {
	Path string `json:"path"`
}

// MultiPartForm represents a multi part form.
type MultiPartForm struct {
	FormData            *io.PipeReader
	FormDataContentType string
}

// GetFileNameFromContentDisposition returns filename from Content-Disposition header.
func GetFileNameFromContentDisposition(header string) (string, error) {
	lowerHeader := strings.ToLower(header)
	if idx := strings.Index(lowerHeader, "filename="); idx != -1 {
		start := idx + len("filename=")
		filename := header[start:]

		// ; after filename
		if idx = strings.Index(filename, ";"); idx != -1 {
			filename = filename[:idx]
		}

		// " " space after filename
		if idx = strings.Index(filename, " "); idx != -1 {
			filename = filename[:idx]
		}

		return strings.TrimSpace(filename), nil
	}

	return "", errors.New("invalid header")
}
