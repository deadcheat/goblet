package main

import (
	"log"
	"os"

	gpf "github.com/deadcheat/goblet/generator/presenter/writer"
	grr "github.com/deadcheat/goblet/generator/repository/regexp"
	guf "github.com/deadcheat/goblet/generator/usecase/file"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "goblet"
	app.Usage = "make a binary contain some assets"
	app.Version = "0.2.0"

	app.Flags = []cli.Flag{
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
	// mount presenter
	presenter := gpf.New(guf.New(grr.New()))
	if err := presenter.Mount(app); err != nil {
		log.Fatal(err)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
