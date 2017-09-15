package main

import (
	"os"

	"github.com/urfave/cli"
)

// NAME : The application name
const NAME string = "ulv-auth"

// USAGE : Application cli description
const USAGE string = "Start the ulv-auth OAuth2.0 / OIDC authentication and identity provider server"

// VERSION : current version number
const VERSION string = "0.1.0"

func main() {
	app := cli.NewApp()

	app.Name = NAME
	app.Usage = USAGE
	app.Version = VERSION

	app.Run(os.Args)
}
