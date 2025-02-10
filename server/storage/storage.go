package storage

import "github.com/fmich7/fyle/types"

type Storage interface {
	UploadFile(file *types.File) error
	DownloadFile(path string) error
}
