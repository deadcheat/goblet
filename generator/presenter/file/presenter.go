package file

import (
	"errors"
	"io"

	"github.com/deadcheat/awsset/generator"
	"github.com/urfave/cli"
)

// Presenter acts for presentation
type Presenter struct {
	w       io.Writer
	usecase generator.UseCase
}

// New presenter
func New(w io.Writer, u generator.UseCase) *Presenter {
	return &Presenter{
		w:       w,
		usecase: u,
	}
}

func (p *Presenter) action(c *cli.Context) error {
	return nil
}

// Mount action
func (p *Presenter) Mount(i interface{}) error {
	c, ok := i.(*cli.App)
	if !ok {
		panic(errors.New("illegal Mount() call"))
	}
	c.Action = p.action
	return nil
}
