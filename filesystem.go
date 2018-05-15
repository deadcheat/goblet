package awsset

import (
	"net/http"
)

// FileSystem exported to a generated asset file
type FileSystem struct {
	Dirs  map[string][]string
	Files map[string]*File
}

// NewFS return newFS FileSystem pointer
func NewFS(dirs map[string][]string, files map[string]*File) *FileSystem {
	return &FileSystem{Dirs: dirs, Files: files}
}

// Exists check whether file exists
func (fs *FileSystem) Exists(name string) bool {
	_, ok := fs.Files[name]
	return ok
}

// Open file from name
func (fs *FileSystem) Open(name string) (http.File, error) {
	f, ok := fs.Files[name]
	if !ok {
		return nil, ErrFileNotFound
	}
	return f, nil
}
