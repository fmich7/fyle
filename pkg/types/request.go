package types

// LoginResponse contains data that is returned to client after login
type LoginResponse struct {
	Token string `json:"token"`
	Salt  string `json:"salt"`
}

// DownloadRequest is contains data that server expects from client
type DownloadRequest struct {
	Path string `json:"path"`
}

// AuthUserRequest contains data that is used to create/login to acc
type AuthUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// ListFilesRequest contains data that is used to list files in user's dir
type ListFilesRequest struct {
	Path string `json:"path"`
}
