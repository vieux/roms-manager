package cmd

import (
	"fmt"
	"sort"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vrgl117-games/roms-manager/gamelist"
)

func view(c *cli.Context, gamelistFiles []*gamelist.File) {
	games := []*gamelist.Game{}
	for _, gamelistFile := range gamelistFiles {
		for i, game := range gamelistFile.Games {
			if (c.Bool("favorite-only") && game.Favorite || !c.Bool("favorite-only")) && (c.Bool("hidden-only") && game.Hidden || c.Bool("visible-only") && !game.Hidden || !c.Bool("hidden-only") && !c.Bool("visible-only")) {
				games = append(games, &gamelistFile.Games[i])
			}
		}
	}

	sort.SliceStable(games, func(i, j int) bool {
		return games[i].Name < games[j].Name
	})

	for _, game := range games {
		fields := log.Fields{"rom": game.RomName, "system": game.System}
		if !c.Bool("favorite-only") {
			fields["favorite"] = game.Favorite
		}
		if !c.Bool("hidden-only") && !c.Bool("visible-only") {
			fields["hidden"] = game.Hidden
		}
		if game.Hidden {
			fields["reason"] = game.Reason
		}
		log.WithFields(fields).Info(game.Name)
	}
}

func NewViewCmd() *cli.Command {
	return &cli.Command{
		Name:  "view",
		Usage: "view games from a gamelist.xml file",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "gamelist",
				Usage:    "path to the <gamelist.xml> file(s)",
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
			&cli.BoolFlag{
				Name:  "favorite-only",
				Usage: "show only favorite games",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("debug") {
				log.SetLevel(log.DebugLevel)
			}

			gamelistFiles := []*gamelist.File{}
			for _, gamelistPath := range c.StringSlice("gamelist") {
				gamelistFile, err := gamelist.New(gamelistPath)
				if err != nil {
					return fmt.Errorf("unable to open: %s %v", gamelistPath, err)
				}

				log.WithFields(log.Fields{"games": len(gamelistFile.Games), "path": gamelistFile.ShortPath}).Info("gamelist loaded")
				gamelistFiles = append(gamelistFiles, gamelistFile)
			}

			view(c, gamelistFiles)

			return nil
		},
	}
}
