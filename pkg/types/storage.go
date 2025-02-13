package types

type Storage interface {
	UploadFile(file *File) error
	DownloadFile(path string) error
	GetFileUploadsLocation() string
}
