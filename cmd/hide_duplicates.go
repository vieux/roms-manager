package cmd

import (
	"errors"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vrgl117-games/roms-manager/gamelist"
)

func hideDuplicates(c *cli.Context, gamelistFiles []*gamelist.File) {
	for i, leftGamelistFile := range gamelistFiles {
		for _, rightGamelistFile := range gamelistFiles[i+1:] {
			total := 0
			hidden := 0
			for j, game := range rightGamelistFile.Games {
				if !rightGamelistFile.Games[j].Hidden {
					total++
					if _, ok := leftGamelistFile.RomNames[game.RomName]; ok {
						log.WithFields(log.Fields{"rom": game.RomName}).Warnf("already present in %s", leftGamelistFile.ShortPath)
						rightGamelistFile.Games[j].Hidden = true
						if c.Bool("override-favorites") {
							rightGamelistFile.Games[j].Favorite = false
						}
						hidden++
					} else if _, ok := leftGamelistFile.Names[game.Name]; ok {
						log.WithFields(log.Fields{"rom": game.RomName}).Warnf("already present in %s", leftGamelistFile.ShortPath)
						rightGamelistFile.Games[j].Hidden = true
						if c.Bool("override-favorites") {
							rightGamelistFile.Games[j].Favorite = false
						}
						hidden++
					}
				}
			}
			log.WithFields(log.Fields{"path": rightGamelistFile.ShortPath}).Infof("%d games out of %d were marked as hidden", hidden, total)
		}
	}
}

func NewHideDuplicatesCmd() *cli.Command {
	return &cli.Command{
		Name:  "hide-duplicates",
		Usage: "hide duplicates games from a list of gamelist.xml files",
		Flags: []cli.Flag{
			&cli.StringSliceFlag{
				Name:     "gamelist",
				Usage:    "path to the <gamelist.xml> files",
				Required: true,
			},
			&cli.BoolFlag{
				Name:  "override-favorites",
				Usage: "unfavorite hidden games",
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

			if len(gamelistFiles) < 2 {
				return errors.New("at least 2 gameslist.xml files are required to hide duplicates")
			}

			hideDuplicates(c, gamelistFiles)

			for _, gamelistFile := range gamelistFiles[1:] {
				if err := gamelistFile.Save(); err != nil {
					return fmt.Errorf("unable to save: %s %v", gamelistFile.Path, err)
				}
			}

			return nil
		},
	}
}
