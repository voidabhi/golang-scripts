package main

import (
	"fmt"
	"os"

	"github.com/codegangsta/cli"
)

func main() {
	app := initApp()
	app.Run(os.Args)
}

func initApp() *cli.App {
	app := cli.NewApp()

	app.Name = "greet"
	app.Version = "0.0.1"
	app.Usage = "Just greet you."
	app.Author = "Tomoya Kose (mitsuse)"
	app.Email = "tomoya@mitsuse.jp"

	app.Commands = []cli.Command{
		newPrintCommand(),
		// Add more sub-commands ...
	}

	return app
}

func newPrintCommand() cli.Command {
	command := cli.Command{
		Name:      "print",
		ShortName: "p",
		Usage:     "Print the greeting",
		Action:    actionPrint,

		Flags: []cli.Flag{
			cli.StringFlag{
				Name:  "greeting,g",
				Value: "Hello",
				Usage: "The greeting to be shown",
			},

			cli.StringFlag{
				Name:  "name,n",
				Value: "world",
				Usage: "The name of person/something to be greeted",
			},

			// Add more command-line options for "print" ...
		},
	}
	return command
}

func actionPrint(ctx *cli.Context) {
	greeting := ctx.String("greeting")
	name := ctx.String("name")

	fmt.Printf("%s, %s.\n", greeting, name)
}
