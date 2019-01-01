package main

import (
	"log"
	"os"

	gpf "github.com/deadcheat/goblet/generator/presenter/file"
	grr "github.com/deadcheat/goblet/generator/repository/regexp"
	guf "github.com/deadcheat/goblet/generator/usecase/file"
	"github.com/deadcheat/goblet/generator/values"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "goblet"
	app.Usage = "make a binary contain some assets"
	app.Version = "1.0.4"

	app.Flags = values.FlagDefs
	// mount presenter
	presenter := gpf.New(guf.New(grr.New()))
	if err := presenter.Mount(app); err != nil {
		log.Fatal(err)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
