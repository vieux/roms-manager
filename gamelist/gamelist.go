package gamelist

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	XMLName xml.Name `xml:"gameList"`
	Games   []Game   `xml:"game"`

	Map       map[string]*Game `xml:"-"`
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

	RomName string `xml:"-"`
}

func New(path string) (*File, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		return nil, err
	}

	var gamelistFile File
	if err := xml.NewDecoder(f).Decode(&gamelistFile); err != nil {
		return nil, err
	}

	gamelistFile.Path = path
	gamelistFile.ShortPath = filepath.Join(filepath.Base(filepath.Dir(path)), filepath.Base(path))

	gamelistFile.Map = make(map[string]*Game, len(gamelistFile.Games))
	for i := range gamelistFile.Games {
		ext := filepath.Ext(gamelistFile.Games[i].Path)
		gamelistFile.Games[i].RomName = strings.TrimSuffix(filepath.Base(gamelistFile.Games[i].Path), ext)
		gamelistFile.Map[gamelistFile.Games[i].RomName] = &gamelistFile.Games[i]
	}

	return &gamelistFile, nil
}

func (gamelistFile *File) Save() error {
	f, err := os.Create(gamelistFile.Path + ".new")
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
