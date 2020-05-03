package main

import (
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vrgl117-games/roms-manager/cmd"
)

func main() {
	log.SetOutput(os.Stdout)

	app := cli.NewApp()
	app.Usage = "a simple tool to scan arcade romset .dat and gamelist.xml files."
	app.Version = "1.1.0"
	app.Authors = []*cli.Author{
		{Name: "Victor Vieux", Email: "github@vrgl117.games"},
	}

	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:  "debug",
			Usage: "enable debug logs",
			Value: false,
		},
	}

	app.Commands = []*cli.Command{
		cmd.NewViewCmd(),
		cmd.NewScanCmd(),
		cmd.NewHideDuplicatesCmd(),
		cmd.NewResetVisibilityCmd(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
