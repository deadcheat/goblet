package dotfileignorematcher

import (
	"testing"

	"github.com/deadcheat/goblet/generator"
)

func TestPrepareSuccess(t *testing.T) {
	e := generator.OptionFlagEntity{
		IgnoreDotFiles: false,
	}

	testee := New()

	err := testee.Prepare(e)
	if err != nil {
		t.Error("Prepare() should not return any errors")
	}
}

func TestMatchReturnTrueWhenIgnoreDotFilesIsFalse(t *testing.T) {
	e := generator.OptionFlagEntity{
		IgnoreDotFiles: false,
	}
	testee := New()
	testee.Prepare(e)

	testPath := ".gitignore"
	if testee.Match(testPath) != true {
		t.Error("Match() should return true when IgnoreDotFiles is false")
	}
}

func TestMatchReturnFalseWhenIgnoreDotFilesIsTrueAndHaveDotPrefix(t *testing.T) {
	e := generator.OptionFlagEntity{
		IgnoreDotFiles: true,
	}
	testee := New()
	testee.Prepare(e)

	testPath := ".gitignore"
	if testee.Match(testPath) != false {
		t.Error("Match() should return false when IgnoreDotFiles is true and has dot-prefix")
	}
}
