package gamelist

import (
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type File struct {
	XMLName xml.Name `xml:"gameList"`
	Games   []Game   `xml:"game"`

	RomNames  map[string]*Game `xml:"-"`
	Names     map[string]*Game `xml:"-"`
	Path      string           `xml:"-"`
	ShortPath string           `xml:"-"`
}

type Game struct {
	Name        string  `xml:"name"`
	Description string  `xml:"desc,omitempty"`
	Path        string  `xml:"path"`
	Publisher   string  `xml:"publisher,omitempty"`
	Developer   string  `xml:"developer,omitempty"`
	ReleaseDate string  `xml:"releasedate,omitempty"`
	Image       string  `xml:"image,omitempty"`
	Video       string  `xml:"video,omitempty"`
	Genre       string  `xml:"genre,omitempty"`
	Rating      float64 `xml:"rating,omitempty"`
	Players     string  `xml:"players,omitempty"`
	Playcount   int     `xml:"playcount,omitempty"`
	LastPlayed  string  `xml:"lastplayed,omitempty"`
	Favorite    bool    `xml:"favorite,omitempty"`
	Hidden      bool    `xml:"hidden,omitempty"`
	Reason      string  `xml:"reason,omitempty"`

	RomName string `xml:"-"`
}

func New(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var gamelistFile File
	if err := xml.NewDecoder(f).Decode(&gamelistFile); err != nil {
		return nil, err
	}

	gamelistFile.Path = path
	gamelistFile.ShortPath = filepath.Join(filepath.Base(filepath.Dir(path)), filepath.Base(path))

	gamelistFile.RomNames = make(map[string]*Game, len(gamelistFile.Games))
	gamelistFile.Names = make(map[string]*Game, len(gamelistFile.Games))

	for i := range gamelistFile.Games {
		ext := filepath.Ext(gamelistFile.Games[i].Path)
		gamelistFile.Games[i].RomName = strings.TrimSuffix(filepath.Base(gamelistFile.Games[i].Path), ext)
		gamelistFile.RomNames[gamelistFile.Games[i].RomName] = &gamelistFile.Games[i]
		gamelistFile.Names[gamelistFile.Games[i].Name] = &gamelistFile.Games[i]

	}

	return &gamelistFile, nil
}

func copyFile(fromPath, toPath string) error {
	from, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return err
}

func (gamelistFile *File) Save() error {
	if err := copyFile(gamelistFile.Path, fmt.Sprintf("%s.old.%d", gamelistFile.Path, time.Now().Unix())); err != nil {
		return err
	}

	f, err := os.Create(gamelistFile.Path)
	if err != nil {
		return err
	}
	defer f.Close()

	content, err := xml.MarshalIndent(gamelistFile, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write([]byte(xml.Header + string(content)))
	return err
}
