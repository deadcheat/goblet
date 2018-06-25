package regexp

import (
	"regexp"

	"github.com/deadcheat/goblet/generator"
)

// Repository implements generator.RegexpRepository
type Repository struct {
	rs []*regexp.Regexp
}

// New return new repository
func New() generator.RegexpRepository {
	return &Repository{}
}

// CompilePatterns compile patterns
func (r *Repository) CompilePatterns(patterns []string) error {
	r.r = make([]*regexp.Regexp, len(patterns))
	for i := range patterns {
		pattern := patterns[i]
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		r.r[i] = reg
	}
	return nil
}

// MatchAny check regexp slices if path matches anyone
func (r *Repository) MatchAny(path string) bool {
	// if no patterns are compiled, return true
	if len(r.rs) == 0 {
		return true
	}
	for i := range r.rs {
		regexp := r.rs[i]
		if regexp != nil && regexp.MatchString(path) {
			return true
		}
	}
	return false
}
