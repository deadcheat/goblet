package file

import (
	"testing"

	"github.com/deadcheat/goblet/generator/mock"
	"github.com/golang/mock/gomock"

	"github.com/urfave/cli"
)

func TestMount(t *testing.T) {
	// Success Pattern
	a := &cli.App{}

	// Prepare mock
	c := gomock.NewController(t)
	defer c.Finish()

	m := mock.NewMockUseCase(c)

	p := New(m)

	if err := p.Mount(a); err != nil {
		t.Error("Mount should not return any errors", err)
	}

}
