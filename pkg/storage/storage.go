package storage

import "github.com/fmich7/fyle/pkg/types"

type Storage interface {
	UploadFile(file *types.File) error
	DownloadFile(path string) error
}
