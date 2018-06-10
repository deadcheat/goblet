package file

import (
	"flag"
	"testing"

	"github.com/deadcheat/goblet/generator/mock"
	"github.com/deadcheat/goblet/generator/values"
	"github.com/golang/mock/gomock"

	"github.com/urfave/cli"
)

func TestMount(t *testing.T) {

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	m := mock.NewMockUseCase(c)
	p := New(m)
	// Success Pattern
	a := &cli.App{}
	if err := p.Mount(a); err != nil {
		t.Error("Mount should not return any errors", err)
	}

	// error pattern
	if err := p.Mount(struct{}{}); err != ErrIllegalMount {
		t.Error("Mount should return ErrIllegalMount", err)
	}
}

var fs = flag.NewFlagSet("test", flag.ContinueOnError)

func TestActionNoArgs(t *testing.T) {

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	m := mock.NewMockUseCase(c)
	p := New(m)
	a := &cli.App{}
	p.Mount(a)
	a.Flags = values.FlagDefs
	ctx := cli.NewContext(a, fs, nil)
	if err := p.action(ctx); err != ErrNoArguments {
		t.Error("Mount should return ErrNoArguments", err)
	}
}
