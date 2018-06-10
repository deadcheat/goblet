package values

import "github.com/urfave/cli"

var FlagDefs = []cli.Flag{
	cli.StringSliceFlag{
		Name:  "expression, e",
		Usage: "Regular expressions you want files to contain",
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
	cli.StringFlag{
		Name:  "name, n",
		Value: "Assets",
		Usage: "Variable name for output assets",
	},
}
