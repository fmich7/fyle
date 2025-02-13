package types

// Storage represents a storage interface for server
type Storage interface {
	UploadFile(file *File) error
	DownloadFile(path string) error
	GetFileUploadsLocation() string
}
