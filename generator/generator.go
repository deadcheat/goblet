package generator

import (
	"github.com/deadcheat/awsset"
)

// Entity generator entity
type Entity struct {
	FsMap map[string]*awsset.FileSystem
	Paths []string
}

// UseCase interface
type UseCase interface {
	LoadFiles(paths []string, ignorePatterns []string) (*Entity, error)
}

// RegexpRepository repository for slice of regexp
type RegexpRepository interface {
	CompilePatterns(patterns []string) error
	MatchAny(path string) bool
}
