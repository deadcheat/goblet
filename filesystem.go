package awsset

import (
	"errors"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/k0kubun/pp"
)

// FileSystem exported to a generated asset file
type FileSystem struct {
	Dirs       map[string][]string
	Files      map[string]*File
	pathPrefix string
}

// NewFS return newFS FileSystem pointer
func NewFS(dirs map[string][]string, files map[string]*File) *FileSystem {
	return &FileSystem{Dirs: dirs, Files: files}
}

// WithPrefix set FileSystem.pathPrefix and return FileSystem itself
func (fs *FileSystem) WithPrefix(prefix string) *FileSystem {
	if fs == nil {
		panic(errors.New("FileSystem.WithPrefix should be called with non-nil receiver"))
	}
	fs.pathPrefix = prefix
	return fs
}

// Exists check whether file exists
func (fs *FileSystem) Exists(name string) bool {
	_, ok := fs.Files[fs.nameResolute(name)]
	return ok
}

// Open file from name
func (fs *FileSystem) Open(name string) (http.File, error) {
	f, ok := fs.Files[fs.nameResolute(name)]
	pp.Println(fs.nameResolute(name), f, ok)
	if !ok {
		return nil, ErrFileNotFound
	}
	return f, nil
}

func (fs *FileSystem) nameResolute(name string) string {
	if name != "" && strings.HasPrefix(name, fs.pathPrefix) {
		pp.Println(name, fs.pathPrefix, filepath.Join("/", strings.TrimPrefix(name, fs.pathPrefix)))
		return filepath.Join("/", strings.TrimPrefix(name, fs.pathPrefix))
	}
	return name
}
