package types

import (
	"bytes"
	"mime/multipart"
	"net/textproto"
)

// File represents a file object
type File struct {
	Filename string
	Header   textproto.MIMEHeader
	Size     int64
	Data     multipart.File
	Owner    string
	Location string
}

// NewFile creates a new File object
func NewFile(header *multipart.FileHeader, file multipart.File, owner string, location string) *File {
	return &File{
		Filename: header.Filename,
		Header:   header.Header,
		Size:     header.Size,
		Data:     file,
		Owner:    owner,
		Location: location,
	}
}

type MultiPartForm struct {
	FormData            *bytes.Buffer
	FormDataContentType string
}
