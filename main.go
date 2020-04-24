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
	app.Version = "0.0.1"
	app.Authors = []*cli.Author{
		{Name: "Victor Vieux", Email: "github@vrgl117.games"},
	}

	app.Commands = []*cli.Command{
		cmd.NewViewCmd(),
		cmd.NewScanCmd(),
		cmd.NewHideDuplicatesCmd(),
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
