package awsset

import (
	"errors"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	if !ok {
		return nil, ErrFileNotFound
	}
	return f, nil
}

func (fs *FileSystem) nameResolute(name string) string {
	if fs.pathPrefix != "" && strings.HasPrefix(name, fs.pathPrefix) {
		return filepath.Join("/", strings.TrimPrefix(name, fs.pathPrefix))
	}
	return name
}

// File returns file struct
func (fs *FileSystem) File(filename string) (*File, error) {
	f, ok := fs.Files[fs.nameResolute(filename)]
	if !ok {
		return nil, ErrFileNotFound
	}
	return f, nil
}

// ReadFile read file and return []byte like as ioutil.ReadFile
func (fs *FileSystem) ReadFile(filename string) ([]byte, error) {
	f, ok := fs.Files[fs.nameResolute(filename)]
	if !ok {
		return nil, ErrFileNotFound
	}
	return f.Data, nil
}

// ReadDir return all files in specified directory
func (fs *FileSystem) ReadDir(dirname string) ([]os.FileInfo, error) {
	dirs, ok := fs.Dirs[dirname]
	if !ok {
		return nil, ErrFileNotFound
	}
	files := make([]os.FileInfo, 0)
	for _, dir := range dirs {
		path := filepath.Join(dirname, dir)
		file, ok := fs.Files[path]
		if !ok {
			continue
		}
		files = append(files, file)
	}
	return files, nil
}
