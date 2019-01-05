package dotfileignorematcher

import "github.com/deadcheat/goblet/generator"

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

// Match return true when path is dotfile path and ignoreDotFiles is true
func (r *Repository) Match(path string) bool {
	return false
}
