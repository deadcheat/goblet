package values

import (
	"bytes"
	"fmt"
	"path/filepath"

	"github.com/urfave/cli"
)

// FlagReader converter flag value to string
type FlagReader func(*cli.Context) string

const (
	// FlagKeyExpression flag name "expression"
	FlagKeyExpression = "expression"
	// FlagKeyGenerate flag name "generate"
	FlagKeyGenerate = "generate"
	// FlagKeyName flag name "name"
	FlagKeyName = "name"
	// FlagKeyOut flag name "out"
	FlagKeyOut = "out"
	// FlagKeyPackage flag name "package"
	FlagKeyPackage = "package"
)

// FlagDefs exported variables for flag defnition
var (
	FlagKeys = []string{
		FlagKeyExpression,
		FlagKeyGenerate,
		FlagKeyName,
		FlagKeyOut,
		FlagKeyPackage,
	}
	FlagReaderMap = map[string]FlagReader{
		FlagKeyExpression: func(c *cli.Context) string {
			es := c.StringSlice(FlagKeyExpression)
			var buf bytes.Buffer
			for _, e := range es {
				_, _ = buf.WriteString(fmt.Sprintf("-%s %s", "e", e))
			}
			return buf.String()
		},
		FlagKeyGenerate: func(c *cli.Context) string {
			if c.Bool(FlagKeyGenerate) {
				return "-g"
			}
			return ""
		},
		FlagKeyName: func(c *cli.Context) string { return fmt.Sprintf("-%s %s", "n", c.String(FlagKeyName)) },
		FlagKeyOut: func(c *cli.Context) string {
			path := c.String(FlagKeyOut)
			bn := filepath.Base(path)
			return fmt.Sprintf("-%s %s", "o", bn)
		},
		FlagKeyPackage: func(c *cli.Context) string { return fmt.Sprintf("-%s %s", "p", c.String(FlagKeyPackage)) },
	}
	FlagDefs = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "expression, e",
			Usage: "Regular expressions you want files to contain",
		},
		cli.BoolFlag{
			Name:  "generate, g",
			Usage: "If set, generate go:generate line to outputfile",
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
