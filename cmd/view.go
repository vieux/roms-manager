package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vrgl117-games/roms-manager/gamelist"
)

func view(c *cli.Context, gamelistFile *gamelist.File) {
	for _, game := range gamelistFile.Games {
		if c.Bool("hidden-only") && game.Hidden || c.Bool("visible-only") && !game.Hidden || !c.Bool("hidden-only") && !c.Bool("visible-only") {
			log.WithFields(log.Fields{"rom": game.RomName, "hidden": game.Hidden}).Info(game.Name)
		}
	}
}

func NewViewCmd() *cli.Command {
	return &cli.Command{
		Name:  "view",
		Usage: "view games from a gamelist.xml file",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "gamelist",
				Usage:    "path to the <gamelist.xml> file",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "hidden-only",
				Usage: "show only hidden games",
			},
			&cli.BoolFlag{
				Name:  "visible-only",
				Usage: "show only visible games",
			},
		},
		Action: func(c *cli.Context) error {
			gamelistFile, err := gamelist.New(c.Path("gamelist"))
			if err != nil {
				log.Fatalf("unable to open: %s %v", c.Path("gamelist"), err)

			}
			log.WithFields(log.Fields{"games": len(gamelistFile.Games), "path": gamelistFile.ShortPath}).Info("gamelist loaded")

			view(c, gamelistFile)

			return nil
		},
	}
}
