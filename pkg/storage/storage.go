package storage

import "github.com/fmich7/fyle/internal/types"

type Storage interface {
	UploadFile(file *types.File) error
	DownloadFile(path string) error
}
