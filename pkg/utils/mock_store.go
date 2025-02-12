package utils

import "github.com/fmich7/fyle/pkg/types"

type MockStore struct{}

func (m *MockStore) UploadFile(file *types.File) error {
	return nil
}

func (m *MockStore) DownloadFile(path string) error {
	return nil
}
