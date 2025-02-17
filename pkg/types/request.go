package types

// DownloadRequest is contains data that server expects from client
type DownloadRequest struct {
	Path string `json:"path"`
}
