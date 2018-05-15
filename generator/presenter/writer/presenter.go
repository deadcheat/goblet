package file

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"os"
	"path/filepath"

	"github.com/deadcheat/awsset/generator"
	pt "github.com/deadcheat/awsset/generator/presenter/writer/template"
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
	e, err := p.usecase.LoadFiles(paths, ignores)
	if err != nil {
		return err
	}
	t, err := template.New("asset").Funcs(
		template.FuncMap{
			"sha1": func(s string) string {
				h := sha1.New()
				h.Write([]byte(s))
				r := h.Sum(nil)
				return fmt.Sprintf("%x", r)
			},
		},
	).Parse(pt.AssetFileTemplate)
	if err != nil {
		return err
	}

	var b bytes.Buffer
	if err = t.Execute(&b, &pt.Assets{
		PackageName: c.String("package"),
		VarName:     c.String("var"),
		DirMap:      e.DirMap,
		FileMap:     e.FileMap,
		Paths:       e.Paths,
	}); err != nil {
		return err
	}

	// gofmt
	formatted, err := format.Source(b.Bytes())

	var writer io.Writer = os.Stdout
	outName := c.String("out")
	if outName != "" {
		// current dir
		target, _ := filepath.Abs(outName)
		f, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer f.Close()
		writer = f
	}
	if _, err = fmt.Fprintln(writer, string(formatted)); err != nil {
		return err
	}
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
