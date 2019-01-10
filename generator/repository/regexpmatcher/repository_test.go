package regexpmatcher

import (
	"testing"

	"github.com/deadcheat/goblet/generator"
)

func TestPrepare(t *testing.T) {
	rep := New()
	// Success pattern
	err := rep.Prepare(generator.OptionFlagEntity{
		IncludePatterns: []string{
			`^[a-z0-9]+\.html$`,
			`^[a-z0-9]+\.js|css$`,
		}})
	if err != nil {
		t.Error("Prepare should not return any errors", err)
	}
	// Fail pattern
	err = rep.Prepare(generator.OptionFlagEntity{
		IncludePatterns: []string{
			`^[a-z0-9]+\.html$`,
			`[a-9]`,
			`^[a-z0-9]+\.js|css$`,
		}})
	if err == nil {
		t.Error("Prepare should  return any error")
	}
}

func TestMatch(t *testing.T) {
	rep := New()
	// No patterns
	if !rep.Match("anything ok") {
		t.Error("if there are not any specified pattern func should return always true")
	}

	// Success pattern
	err := rep.Prepare(generator.OptionFlagEntity{
		IncludePatterns: []string{
			`^[a-z0-9]+\.html$`,
			`^[a-z0-9]+\.js|css$`,
		}})
	if err != nil {
		t.Error("Prepare should not return any errors", err)
	}

	// True pattern
	pat := "hogehoge.html"
	if !rep.Match(pat) {
		t.Error(pat, " should be matched but not")
	}
	// False pattern
	pat = "thisisnotmatched.txt"
	if rep.Match(pat) {
		t.Error(pat, " should not be matched but matched")
	}
}
