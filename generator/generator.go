package generator

import (
	"github.com/deadcheat/goblet"
)

// Entity generator entity
type Entity struct {
	DirMap  map[string][]string
	FileMap map[string]*goblet.File
	Paths   []string
}

// UseCase interface
type UseCase interface {
	LoadFiles(paths []string, includePatterns []string) (*Entity, error)
}

// RegexpRepository repository for slice of regexp
type RegexpRepository interface {
	CompilePatterns(patterns []string) error
	MatchAny(path string) bool
}
