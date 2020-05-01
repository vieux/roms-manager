package cmd

import (
	"archive/zip"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"github.com/vrgl117-games/roms-manager/dat"
	"github.com/vrgl117-games/roms-manager/gamelist"
)

func scanGame(c *cli.Context, datGame *dat.Game, game *gamelist.Game) string {
	if strings.Contains(game.Name, "notgame") {
		return "not a game"
	}
	if datGame.Video.Orientation != "" && c.Bool("hide-horizontal") && datGame.Video.Orientation == "horizontal" {
		return "horizontal game"
	}

	if datGame.Video.Orientation != "" && c.Bool("hide-vertical") && datGame.Video.Orientation == "vertical" {
		return "vertical game"
	}

	if aspectRatio := datGame.Video.AspectRatio(); aspectRatio != "x" && c.String("aspect-ratio") != aspectRatio {
		return fmt.Sprintf("incompatible aspect ratio %q", aspectRatio)
	}

	if datGame.Input.Buttons != 0 && datGame.Input.Buttons > c.Int("max-buttons") {
		return "game require too many buttons"
	}

	if datGame.Input.Control != "" {
		if _, ok := c.Generic("controls").(*MapFlag).values[datGame.Input.Control]; !ok {
			return fmt.Sprintf("unknown %q control", datGame.Input.Control)
		}
	}
	for _, hide := range c.StringSlice("forbidden-keywords") {
		if strings.Contains(strings.ToLower(datGame.Description), hide) {
			return fmt.Sprintf("%q found in description", hide)
		}

		if datGame.Manufacturer == hide {
			return fmt.Sprintf("manufacturer is %q", hide)
		}
	}

	return ""
}

func scanZip(datGame *dat.Game, game *gamelist.Game, path string) string {
	r, err := zip.OpenReader(path)
	if err != nil {
		return fmt.Sprintf("unable to open: %s %v", path, err)
	}
	defer r.Close()

	for _, romFile := range r.File {
		if rom, ok := datGame.RomNames[romFile.Name]; ok {
			if rom.Size != int(romFile.FileHeader.UncompressedSize) {
				return fmt.Sprintf("file %s has wrong size", romFile.Name)
			}

			crc, err := strconv.ParseUint(rom.CRC, 16, 32)
			if err != nil {
				return fmt.Sprintf("unable to convert CRC: %s %v", rom.CRC, err)
			}
			if uint32(crc) != romFile.FileHeader.CRC32 {
				return fmt.Sprintf("file %s has wrong CRC", romFile.Name)
			}
		}
	}

	return ""
}

func hideGame(c *cli.Context, gamelistFile *gamelist.File, i int, reason string) {
	log.WithFields(log.Fields{"rom": gamelistFile.Games[i].RomName}).Warn(reason)
	gamelistFile.Games[i].Hidden = true
	gamelistFile.Games[i].Reason = reason
	if c.Bool("override-favorites") {
		gamelistFile.Games[i].Favorite = false
	}
}

func scanGames(c *cli.Context, datFile *dat.File, gamelistFile *gamelist.File) (int, error) {
	hidden := 0

	for i, game := range gamelistFile.Games {
		log.WithFields(log.Fields{"rom": game.RomName, "hidden": game.Hidden}).Debugf("scanning game")

		if c.Bool("force-visible") {
			gamelistFile.Games[i].Hidden = false
		}
		if game.Hidden {
			continue
		}
		if datGame, ok := datFile.RomNames[game.RomName]; !ok {
			hideGame(c, gamelistFile, i, "not present in .dat file")
			hidden++
		} else if reason := scanGame(c, datGame, &game); reason != "" {
			hideGame(c, gamelistFile, i, reason)
			hidden++
		} else if c.Bool("zip") {
			if reason := scanZip(datGame, &game, filepath.Join(filepath.Dir(gamelistFile.Path), game.Path)); reason != "" {
				hideGame(c, gamelistFile, i, reason)
				hidden++
			}
		}
	}

	return hidden, nil
}

func scanClones(c *cli.Context, datFile *dat.File, gamelistFile *gamelist.File) (int, error) {
	hidden := 0
	for _, game := range gamelistFile.Games {
		if game.Hidden {
			continue
		}
		datGame, ok := datFile.RomNames[game.RomName]
		if !ok {
			continue
		}
		if datGame.CloneOf == "" {
			continue
		}
		log.WithFields(log.Fields{"rom": game.RomName, "hidden": game.Hidden}).Debugf("scanning clone")

		if parentDatGame, ok := datFile.RomNames[datGame.CloneOf]; ok {
			visibleGamelistClones := []*gamelist.Game{}

			visibleGamelistClones = append(visibleGamelistClones, gamelistFile.RomNames[datGame.CloneOf])
			for _, datClone := range parentDatGame.Clones {
				if gamelistClone, ok := gamelistFile.RomNames[datClone.Name]; ok && !gamelistClone.Hidden {
					visibleGamelistClones = append(visibleGamelistClones, gamelistClone)
				}
			}

			if len(visibleGamelistClones) > 1 {
				sort.SliceStable(visibleGamelistClones, func(i, j int) bool {
					return datFile.RomNames[visibleGamelistClones[i].RomName].Description > datFile.RomNames[visibleGamelistClones[j].RomName].Description
				})

				num := 0
				if c.Bool("clones-selection") {
					fmt.Printf("[%s] %s has %d clones:\n", parentDatGame.Name, parentDatGame.Description, len(visibleGamelistClones)-1)
					for i, visibleGamelistClone := range visibleGamelistClones {
						fmt.Printf(" %d: [%s] %s\n", i, visibleGamelistClone.RomName, datFile.RomNames[visibleGamelistClone.RomName].Description)
					}

					fmt.Printf("Keep game #: ")
					fmt.Scanf("%d\n", &num)
				}

				for i := range visibleGamelistClones {
					if i != num {
						visibleGamelistClones[i].Hidden = true
						if c.Bool("override-favorites") {
							visibleGamelistClones[i].Favorite = false
						}
						hidden++
					}
				}
			}

		}
	}

	return hidden, nil
}

func NewScanCmd() *cli.Command {
	return &cli.Command{
		Name:  "scan",
		Usage: "scan a gamelist.xml file and hide incompatible games",
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:     "dat",
				Usage:    "path to the <.dat> or <.xml> file",
				Required: true,
			},
			&cli.PathFlag{
				Name:     "gamelist",
				Usage:    "path to the <gamelist.xml> file",
				Required: true,
			},
			&cli.StringSliceFlag{
				Name:  "forbidden-keywords",
				Usage: "hide games if the following keywords are found in the description",
				Value: cli.NewStringSlice("hack", "bootleg", "homebrew", "prototype", "jamma pcb"),
			},
			&cli.BoolFlag{
				Name:  "clones-selection",
				Usage: "user selects which clone (or original) to show, otherwise, try to automatically pick the best one",
			},
			&cli.BoolFlag{
				Name:  "hide-horizontal",
				Usage: "hide horizontal games",
			},
			&cli.BoolFlag{
				Name:  "hide-vertical",
				Usage: "hide vertical games",
			},
			&cli.BoolFlag{
				Name:  "override-favorites",
				Usage: "unfavorite hidden games",
			},
			&cli.IntFlag{
				Name:  "max-buttons",
				Usage: "hide games that require too many buttons",
				Value: 6,
			},
			&cli.GenericFlag{
				Name:  "controls",
				Usage: "hide games with incompatible controls",
				Value: NewMapFlag("joy8way", "joy4way", "stick"),
			},
			&cli.StringFlag{
				Name:  "aspect-ratio",
				Usage: "hide games with incompatible aspect ratio",
				Value: "4x3",
			},
			&cli.BoolFlag{
				Name:  "zip",
				Usage: "scan the inside of the rom",
				Value: true,
			},
			&cli.BoolFlag{
				Name:  "force-visible",
				Usage: "set the game to visible even if it was hidden if roms-manager determines it is compatible",
			},
		},
		Action: func(c *cli.Context) error {
			if c.Bool("debug") {
				log.SetLevel(log.DebugLevel)
			}

			datFile, err := dat.New(c.Path("dat"))
			if err != nil {
				return fmt.Errorf("unable to open: %s %v", c.Path("dat"), err)

			}
			log.WithFields(log.Fields{"games": len(datFile.Games), "path": datFile.ShortPath}).Info("database loaded")

			gamelistFile, err := gamelist.New(c.Path("gamelist"))
			if err != nil {
				return fmt.Errorf("unable to open: %s %v", c.Path("gamelist"), err)

			}
			log.WithFields(log.Fields{"games": len(gamelistFile.Games), "path": gamelistFile.ShortPath}).Info("gamelist loaded")

			hiddenGames, err := scanGames(c, datFile, gamelistFile)
			if err != nil {
				return err
			}

			hiddenClones, err := scanClones(c, datFile, gamelistFile)
			if err != nil {
				return err
			}

			log.WithFields(log.Fields{"path": gamelistFile.ShortPath}).Infof("%d games out of %d were marked as hidden", hiddenGames+hiddenClones, len(gamelistFile.Games))

			if err := gamelistFile.Save(); err != nil {
				return fmt.Errorf("unable to save: %s %v", gamelistFile.Path, err)
			}
			return nil
		},
	}
}

type MapFlag struct {
	values map[string]struct{}
}

func NewMapFlag(values ...string) *MapFlag {
	mf := MapFlag{values: make(map[string]struct{}, len(values))}
	for _, value := range values {
		_ = mf.Set(value)
	}
	return &mf
}

func (mf *MapFlag) Set(value string) error {
	mf.values[value] = struct{}{}
	return nil
}

func (mf *MapFlag) String() string {
	return fmt.Sprint(mf.values)
}
