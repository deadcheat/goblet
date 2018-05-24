package goblet

import (
	"bytes"
	"os"
	"path"
	"time"
)

// File file system model for each files
type File struct {
	Path       string
	Data       []byte
	FileMode   os.FileMode
	ModifiedAt time.Time
	// implementation for io.Reader and io.Seeker
	buf *bytes.Reader
}

// NewFromFileInfo return new file
func NewFromFileInfo(fi os.FileInfo, p string, d []byte) *File {
	return &File{
		Path:       p,
		Data:       d,
		FileMode:   fi.Mode(),
		ModifiedAt: fi.ModTime(),
	}
}

// NewFile return new file
func NewFile(p string, d []byte, m os.FileMode, mat time.Time) *File {
	return &File{
		Path:       p,
		Data:       d,
		FileMode:   m,
		ModifiedAt: mat,
	}
}

// Name implement for os.FileInfo
func (f *File) Name() string {
	return path.Base(f.Path)
}

// Size implement for os.FileInfo
func (f *File) Size() int64 {
	return int64(len(f.Data))
}

// Mode implement for os.FileInfo
func (f *File) Mode() os.FileMode {
	return f.FileMode
}

// ModTime implement for os.FileInfo
func (f *File) ModTime() time.Time {
	return f.ModifiedAt
}

// IsDir implement for os.FileInfo
func (f *File) IsDir() bool {
	return f.FileMode.IsDir()
}

// Sys implement for os.FileInfo
func (f *File) Sys() interface{} {
	return nil
}

// Read for io.Reader
func (f *File) Read(data []byte) (int, error) {
	if f.buf == nil {
		f.buf = bytes.NewReader(f.Data)
	}

	return f.buf.Read(data)
}

// Seek for io.Seeker
func (f *File) Seek(offset int64, whence int) (int64, error) {
	if f.buf == nil {
		f.buf = bytes.NewReader(f.Data)
	}

	return f.buf.Seek(offset, whence)
}

// Close implement for io.Closer
func (f *File) Close() error {
	f.buf = nil
	return nil
}

// Readdir implement for http.File
func (f *File) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}

// Stat implement for http.File
func (f *File) Stat() (os.FileInfo, error) {
	return f, nil
}
