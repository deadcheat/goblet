package awsset

import (
	"os"
)

// FileSystem file system model for each files
type FileSystem struct {
	FileInfo os.FileInfo
	Path     string
	Data     []byte
}
