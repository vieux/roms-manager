package cmd

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vrgl117-games/roms-manager/gamelist"
)

func hideDuplicates(gamelistFiles []*gamelist.File) {
	for i, leftGamelistFile := range gamelistFiles {
		for _, rightGamelistFile := range gamelistFiles[i+1:] {
			total := 0
			hidden := 0
			for j, game := range rightGamelistFile.Games {
				if !rightGamelistFile.Games[j].Hidden {
					total++
					if _, ok := leftGamelistFile.Map[game.RomName]; ok {
						log.WithFields(log.Fields{"rom": game.RomName}).Warnf("already present in %s", leftGamelistFile.ShortPath)
						rightGamelistFile.Games[j].Hidden = true
						rightGamelistFile.Games[j].Favorite = false
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
				Usage:    "path to the master <gamelist.xml> file(s)",
				Required: true,
			},
		},
		Action: func(c *cli.Context) error {
			gamelistFiles := []*gamelist.File{}
			for _, gamelistPath := range c.StringSlice("gamelist") {
				gamelistFile, err := gamelist.New(gamelistPath)
				if err != nil {
					log.Fatalf("unable to open: %s %v", gamelistPath, err)
					return err
				}

				log.WithFields(log.Fields{"games": len(gamelistFile.Games), "path": gamelistFile.ShortPath}).Info("gamelist loaded")
				gamelistFiles = append(gamelistFiles, gamelistFile)
			}

			hideDuplicates(gamelistFiles)

			for _, gamelistFile := range gamelistFiles[1:] {
				if err := gamelistFile.Save(); err != nil {
					log.Errorf("unable to save: %s %v", gamelistFile.Path, err)
					return err
				}
			}

			return nil
		},
	}
}
