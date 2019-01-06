package dotfileignorematcher

import (
	"log"
	"path/filepath"
	"strings"

	"github.com/deadcheat/goblet/generator"
)

// Repository implements generator.RegexpRepository
type Repository struct {
	ignoreDotFiles bool
}

// New return new repository
func New() generator.PathMatcherRepository {
	return &Repository{}
}

// Prepare set ignoreDotFiles from flag
func (r *Repository) Prepare(e generator.OptionFlagEntiry) error {
	r.ignoreDotFiles = e.IgnoreDotFiles
	return nil
}

// Match return true if path IS NOT A DOTFILE when ignoreDotFiles is true, or else, return true
func (r *Repository) Match(path string) bool {
	if r.ignoreDotFiles {
		abpath, err := filepath.Abs(path)
		if err != nil {
			log.Println("couldn't get abs from path ", path)
			return false
		}
		basepath := filepath.Base(abpath)
		return !strings.HasPrefix(basepath, ".")
	}
	return true
}
