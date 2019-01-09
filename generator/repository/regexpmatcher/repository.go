package regexpmatcher

import (
	"regexp"

	"github.com/deadcheat/goblet/generator"
)

// Repository implements generator.RegexpRepository
type Repository struct {
	rs []*regexp.Regexp
}

// New return new repository
func New() generator.PathMatcherRepository {
	return &Repository{}
}

// CompilePatterns compile patterns
func (r *Repository) Prepare(e generator.OptionFlagEntity) error {
	r.rs = make([]*regexp.Regexp, len(e.IncludePatterns))
	for i := range e.IncludePatterns {
		pattern := e.IncludePatterns[i]
		reg, err := regexp.Compile(pattern)
		if err != nil {
			return err
		}
		r.rs[i] = reg
	}
	return nil
}

// MatchAny check regexp slices if path matches anyone
func (r *Repository) Match(path string) bool {
	// if no patterns are compiled, return true
	if len(r.rs) == 0 {
		return true
	}
	for i := range r.rs {
		reg := r.rs[i]
		if reg != nil && reg.MatchString(path) {
			return true
		}
	}
	return false
}
