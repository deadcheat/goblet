package values

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/urfave/cli"
)

// FlagReader converter flag value to string
type FlagReader func(*cli.Context) string

// FlagDefs exported variables for flag defnition
var (
	FlagKeys = []string{
		"expression",
		"name",
		"out",
		"package",
	}
	FlagReaderMap = map[string]FlagReader{
		"expression": func(c *cli.Context) string {
			es := c.StringSlice("expression")
			var buf bytes.Buffer
			for _, e := range es {
				_, _ = buf.WriteString(fmt.Sprintf("-%s %s", "e", e))
			}
			return buf.String()
		},
		"out": func(c *cli.Context) string {
			path := c.String("out")
			bn := filepath.Base(path)
			return fmt.Sprintf("-%s %s", "o", bn)
		},
		"package": func(c *cli.Context) string { return fmt.Sprintf("-%s %s", "p", c.String("package")) },
		"name":    func(c *cli.Context) string { return fmt.Sprintf("-%s %s", "n", c.String("name")) },
	}
	FlagDefs = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "expression, e",
			Usage: "Regular expressions you want files to contain",
		},
		cli.StringFlag{
			Name:  "name, n",
			Value: "Assets",
			Usage: "Variable name for output assets",
		},
		cli.StringFlag{
			Name:  "out, o",
			Usage: "Output file name, result will be displaed to standard-out when it's skipped",
		},
		cli.StringFlag{
			Name:  "package, p",
			Value: "main",
			Usage: "Package name for output",
		},
	}
)
