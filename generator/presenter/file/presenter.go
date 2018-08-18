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
	"strings"

	"github.com/deadcheat/goblet/generator"
	pt "github.com/deadcheat/goblet/generator/presenter/file/template"
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

var (
	// ErrNoArguments will be returned when didn't specify any arguments
	ErrNoArguments = errors.New("Please specify the argument")
	// ErrIllegalMount will be returned when wrong argument passed to Mount
	ErrIllegalMount = errors.New("illegal Mount() call")
)

func (p *Presenter) action(c *cli.Context) error {
	if c.NArg() == 0 {
		return ErrNoArguments
	}
	paths := append([]string{c.Args().First()}, c.Args().Tail()...)
	ignores := c.StringSlice("except")
	e, err := p.usecase.LoadFiles(paths, ignores)
	if err != nil {
		return err
	}
	t, _ := template.New("asset").Funcs(
		template.FuncMap{
			"sha1": func(s string) string {
				h := sha1.New()
				h.Write([]byte(s))
				r := h.Sum(nil)
				return fmt.Sprintf("%x", r)
			},
			"printData": func(b []byte) template.HTML {
				s := fmt.Sprintf("%#v", string(b))
				return template.HTML(fmt.Sprint(s))
			},
		},
	).Parse(pt.AssetFileTemplate)

	var b bytes.Buffer
	_ = t.Execute(&b, &pt.Assets{
		ExecutedCommand: strings.Join(os.Args, " "),
		PackageName:     c.String("package"),
		VarName:         c.String("name"),
		DirMap:          e.DirMap,
		FileMap:         e.FileMap,
		Paths:           e.Paths,
	})

	// gofmt
	formatted, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}

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
	fmt.Fprintln(writer, string(formatted))
	return nil
}

// Mount action
func (p *Presenter) Mount(i interface{}) error {
	c, ok := i.(*cli.App)
	if !ok {
		return ErrIllegalMount
	}
	c.Action = p.action
	return nil
}
