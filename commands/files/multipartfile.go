package files

import (
	"mime"
	"mime/multipart"
	"net/http"
	"net/url"
)

const (
	multipartFormdataType = "multipart/form-data"
	multipartMixedType    = "multipart/mixed"

	contentTypeHeader = "Content-Type"
)

// MultipartFile implements File, and is created from a `multipart.Part`.
// It can be either a directory or file (checked by calling `IsDirectory()`).
type MultipartFile struct {
	File

	Part      *multipart.Part
	Reader    *multipart.Reader
	Mediatype string
}

func NewFileFromPart(part *multipart.Part) (File, error) {
	f := &MultipartFile{
		Part: part,
	}

	contentType := part.Header.Get(contentTypeHeader)

	var params map[string]string
	var err error
	f.Mediatype, params, err = mime.ParseMediaType(contentType)
	if err != nil {
		return nil, err
	}

	if f.IsDirectory() {
		boundary, found := params["boundary"]
		if !found {
			return nil, http.ErrMissingBoundary
		}

		f.Reader = multipart.NewReader(part, boundary)
	}

	return f, nil
}

func (f *MultipartFile) IsDirectory() bool {
	return f.Mediatype == multipartFormdataType || f.Mediatype == multipartMixedType
}

func (f *MultipartFile) NextFile() (File, error) {
	if !f.IsDirectory() {
		return nil, ErrNotDirectory
	}

	part, err := f.Reader.NextPart()
	if err != nil {
		return nil, err
	}

	return NewFileFromPart(part)
}

func (f *MultipartFile) FileName() string {
	if f == nil || f.Part == nil {
		return ""
	}

	filename, err := url.QueryUnescape(f.Part.FileName())
	if err != nil {
		// if there is a unescape error, just treat the name as unescaped
		return f.Part.FileName()
	}
	return filename
}

func (f *MultipartFile) Read(p []byte) (int, error) {
	if f.IsDirectory() {
		return 0, ErrNotReader
	}
	return f.Part.Read(p)
}

func (f *MultipartFile) Close() error {
	if f.IsDirectory() {
		return ErrNotReader
	}
	return f.Part.Close()
}
