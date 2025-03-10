package types

// DownloadRequest is contains data that server expects from client
type DownloadRequest struct {
	Path string `json:"path"`
}

// CreateUserRequest contains data that is used to create a new acc
type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginUserRequest contains data that is used to create a new acc
type LoginUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
