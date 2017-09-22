package main

import (
	"os"

	"github.com/urfave/cli"

	"github.com/myth/ulv-auth/oidc"
)

// NAME : The application name
const NAME string = "ulv-auth"

// USAGE : Application cli description
const USAGE string = "OAuth2.0 / OIDC authentication and identity provider server"

// VERSION : current version number
const VERSION string = "0.1.0"

func main() {
	app := cli.NewApp()

	app.Name = NAME
	app.Usage = USAGE
	app.Version = VERSION

	app.Commands = []cli.Command{
		{
			Name:    "server",
			Aliases: []string{"s"},
			Usage:   "Starts the Ulv Auth server",
			Action: func(c *cli.Context) error {
				oidc.StartExampleServer()
				return nil
			},
		},
	}

	app.Run(os.Args)
}
