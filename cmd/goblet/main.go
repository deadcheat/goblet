package main

import (
	"log"
	"os"

	"github.com/deadcheat/goblet/generator"
	gpf "github.com/deadcheat/goblet/generator/presenter/file"
	grd "github.com/deadcheat/goblet/generator/repository/dotfileignorematcher"
	grr "github.com/deadcheat/goblet/generator/repository/regexpmatcher"
	guf "github.com/deadcheat/goblet/generator/usecase/file"
	"github.com/deadcheat/goblet/generator/values"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "goblet"
	app.Usage = "make a binary contain some assets"
	app.Version = "1.2.2"

	app.Flags = values.FlagDefs
	// mount presenter
	presenter := gpf.New(guf.New([]generator.PathMatcherRepository{grr.New(), grd.New()}))
	if err := presenter.Mount(app); err != nil {
		log.Fatal(err)
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
