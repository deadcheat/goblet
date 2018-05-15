package regexp

import (
	"regexp"

	"github.com/deadcheat/awsset/generator"
)

type Repository struct {
	r []*regexp.Regexp
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
	for i := range r.r {
		regexp := r.r[i]
		if regexp != nil && regexp.MatchString(path) {
			return true
		}
	}
	return false
}
