package regexp

import "testing"

func TestCompilePatterns(t *testing.T) {
	rep := New()
	// Success pattern
	err := rep.CompilePatterns([]string{
		`^[a-z0-9]+\.html$`,
		`^[a-z0-9]+\.js|css$`,
	},
	)
	if err != nil {
		t.Error("CompilePatterns should not return any errors", err)
	}
	// Fail pattern
	err = rep.CompilePatterns([]string{
		`^[a-z0-9]+\.html$`,
		`[a-9]`,
		`^[a-z0-9]+\.js|css$`,
	},
	)
	if err == nil {
		t.Error("CompilePatterns should  return any error")
	}
}

func TestMatchAny(t *testing.T) {
	rep := New()
	// No patterns
	if !rep.MatchAny("anything ok") {
		t.Error("if there are not any specified pattern func should return always true")
	}

	// Success pattern
	err := rep.CompilePatterns([]string{
		`^[a-z0-9]+\.html$`,
		`^[a-z0-9]+\.js|css$`,
	},
	)
	if err != nil {
		t.Error("CompilePatterns should not return any errors", err)
	}

	// True pattern
	pat := "hogehoge.html"
	if !rep.MatchAny(pat) {
		t.Error(pat, " should be matched but not")
	}
	// False pattern
	pat = "thisisnotmatched.txt"
	if rep.MatchAny(pat) {
		t.Error(pat, " should not be matched but matched")
	}
}
