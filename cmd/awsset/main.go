package main

import (
	"log"
	"os"

	gpf "github.com/deadcheat/awsset/generator/presenter/writer"
	grr "github.com/deadcheat/awsset/generator/repository/regexp"
	guf "github.com/deadcheat/awsset/generator/usecase/file"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "awsset"
	app.Usage = "make a binary contain some assets"

	app.Flags = []cli.Flag{
		cli.StringSliceFlag{
			Name:  "except, e",
			Usage: "Regular expressions you want files to ignore",
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
			Name:  "var, t",
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
