package goblet

import "errors"

var (
	// ErrFileNotFound error when file is not found
	ErrFileNotFound = errors.New("specified file is not found")
)
