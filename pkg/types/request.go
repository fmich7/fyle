package types

// DownloadRequest is contains data that server expects from client
type DownloadRequest struct {
	Path string `json:"path"`
}

// AuthUserRequest contains data that is used to create/login to acc
type AuthUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
