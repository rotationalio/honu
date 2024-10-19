package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/rotationalio/honu/pkg"
	"github.com/rotationalio/honu/pkg/config"
	"github.com/urfave/cli/v2"
)

func main() {
	// If a dotenv file exists load it for configuration
	godotenv.Load()

	// Create HonuDB command line application
	app := cli.NewApp()
	app.Name = "honudb"
	app.Version = pkg.Version()
	app.Usage = "run and manage a honudb replica service"
	app.Flags = []cli.Flag{}
	app.Commands = []*cli.Command{
		{
			Name:     "serve",
			Usage:    "run the honu replica service",
			Category: "server",
			Action:   serve,
			Flags:    []cli.Flag{},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

//===========================================================================
// Server Commands
//===========================================================================

func serve(c *cli.Context) (err error) {
	// Load the configuration from a file or from the environment.
	var conf config.Config
	if conf, err = config.New(); err != nil {
		return cli.Exit(err, 1)
	}

	fmt.Println(conf)
	return nil
}
