package file

import (
	"errors"
	"fmt"

	"github.com/deadcheat/awsset/generator"
	"github.com/urfave/cli"
)

// Presenter acts for presentation
type Presenter struct {
	usecase generator.UseCase
}

// New presenter
func New(u generator.UseCase) *Presenter {
	return &Presenter{
		usecase: u,
	}
}

func (p *Presenter) action(c *cli.Context) error {
	if c.NArg() == 0 {
		return errors.New("Please specify the argument")
	}
	paths := append([]string{c.Args().First()}, c.Args().Tail()...)
	ignores := c.StringSlice("except")
	entities, err := p.usecase.LoadFiles(paths, ignores)
	if err != nil {
		return err
	}
	fmt.Printf("%#v", entities)
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
