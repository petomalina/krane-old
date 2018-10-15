package main

import (
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name = "Krane"
	app.Usage = "manipulates the krane-operator"

	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
}
