package file

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/textproto"
	"os"

	"github.com/spf13/afero"
)

// File represents a file object.
type File struct {
	Filename string
	Header   textproto.MIMEHeader
	Size     int64
	Data     multipart.File
	Owner    string
	Location string
}

// NewFile creates a new File object.
func NewFile(
	header *multipart.FileHeader,
	file multipart.File,
	owner string,
	location string,
) *File {
	return &File{
		Filename: header.Filename,
		Header:   header.Header,
		Size:     header.Size,
		Data:     file,
		Owner:    owner,
		Location: location,
	}
}

// SaveFileOnDisk saves file on disk given its path and content.
// It will not overwrite existing files on your disk.
func SaveFileOnDisk(fs afero.Fs, filepath string, content io.Reader) error {
	// ensure the directory exists
	if err := fs.MkdirAll(filepath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %w", filepath, err)
	}

	// create the new file
	file, err := fs.Create(filepath)
	if err != nil {
		return fmt.Errorf("failed to create file %s: %v", filepath, err)
	}
	defer file.Close()

	// copy the content to the file
	_, err = io.Copy(file, content)
	if err != nil {
		return errors.New("copying data to file")
	}

	return nil
}
