package types

// DownloadRequest is contains data that server expects from client
type DownloadRequest struct {
	Path string `json:"path"`
	// TEMP SOLUTION BEFORE AUTH
	User string `json:"user"`
}
