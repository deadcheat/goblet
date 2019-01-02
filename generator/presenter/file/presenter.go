package file

import (
	"bytes"
	"crypto/sha1"
	"errors"
	"fmt"
	"go/format"
	"html/template"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/deadcheat/goblet/generator"
	pt "github.com/deadcheat/goblet/generator/presenter/file/template"
	"github.com/deadcheat/goblet/generator/values"
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
	ignores := c.StringSlice(values.FlagKeyExpression)
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
	generateGoGen := c.Bool(values.FlagKeyGenerate)
	var b bytes.Buffer
	assets := &pt.Assets{
		ExecutedCommand:    strings.Join(os.Args, " "),
		PackageName:        c.String(values.FlagKeyPackage),
		GenerateGoGenerate: generateGoGen,
		VarName:            c.String(values.FlagKeyName),
		DirMap:             e.DirMap,
		FileMap:            e.FileMap,
		Paths:              e.Paths,
	}
	targetPaths := paths
	var writer io.Writer = os.Stdout
	outName := c.String(values.FlagKeyOut)
	if outName != "" {
		// current dir
		target, _ := filepath.Abs(outName)
		f, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer f.Close()
		writer = f

		if generateGoGen {
			targetPaths = make([]string, len(paths))
			baseDir := filepath.Dir(outName)
			for i, p := range paths {
				rp, err := filepath.Rel(baseDir, p)
				if err != nil {
					log.Println(err)
					targetPaths[i] = paths[i]
					continue
				}
				targetPaths[i] = rp
			}
		}
	}
	if generateGoGen {
		assets.ExecutedCommand = executedCommand(c, targetPaths)
	}
	_ = t.Execute(&b, assets)

	// gofmt
	formatted, err := format.Source(b.Bytes())
	if err != nil {
		return err
	}
	fmt.Fprintln(writer, string(formatted))
	return nil
}

func executedCommand(c *cli.Context, argPaths []string) string {
	var buf bytes.Buffer
	// write command itself
	buf.WriteString(os.Args[0])
	sort.Strings(values.FlagKeys)
	for _, k := range values.FlagKeys {
		if c.IsSet(k) {
			f := values.FlagReaderMap[k]
			buf.WriteRune(' ')
			buf.WriteString(f(c))
		}
	}
	for _, ap := range argPaths {
		buf.WriteRune(' ')
		buf.WriteString(ap)
	}
	return buf.String()
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
