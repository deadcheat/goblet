package main

import (
	"log"
	"os"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()

	app.Name = "awsset"
	app.Usage = "make a binary contain some assets"

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
