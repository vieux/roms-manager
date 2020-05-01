package dat

import (
	"encoding/xml"
	"os"
	"path/filepath"
)

type File struct {
	Header Header `xml:"header"`
	Games  []Game `xml:"game"`

	RomNames  map[string]*Game `xml:"-"`
	Path      string           `xml:"-"`
	ShortPath string           `xml:"-"`
}

type Header struct {
	Name        string `xml:"name"`
	Description string `xml:"description"`
	Category    string `xml:"category"`
	Version     string `xml:"version"`
	Author      string `xml:"author"`
	Homepage    string `xml:"homepage"`
	URL         string `xml:"url"`
}

type Game struct {
	Name         string `xml:"name,attr"`
	RomOf        string `xml:"romof,attr"`
	CloneOf      string `xml:"cloneof,attr"`
	Description  string `xml:"description"`
	Year         string `xml:"year"`
	Manufacturer string `xml:"manufacturer"`
	Video        Video  `xml:"video"`
	Input        struct {
		Buttons int    `xml:"buttons,attr"`
		Control string `xml:"control,attr"`
	} `xml:"input"`
	Roms   []Rom `xml:"rom"`
	Driver struct {
		Status string `xml:"status,attr"`
	} `xml:"driver"`

	RomNames map[string]*Rom  `xml:"-"`
	Clones   map[string]*Game `xml:"-"`
}

type Video struct {
	Orientation string `xml:"orientation,attr"`
	AspectX     string `xml:"aspectx,attr"`
	AspectY     string `xml:"aspecty,attr"`
}

func (v *Video) AspectRatio() string {
	return v.AspectX + "x" + v.AspectY
}

type Rom struct {
	Name    string `xml:"name,attr"`
	Size    int    `xml:"size,attr"`
	Merge   string `xml:"merge,attr"`
	CRC     string `xml:"crc,attr"`
	SHA1    string `xml:"sha1,attr"`
	Region  string `xml:"region,attr"`
	Dispose string `xml:"dispose,attr"`
	Offset  string `xml:"offset,attr"`
}

func New(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var datFile File
	if err := xml.NewDecoder(f).Decode(&datFile); err != nil {
		return nil, err
	}

	datFile.Path = path
	datFile.ShortPath = filepath.Join(filepath.Base(filepath.Dir(path)), filepath.Base(path))

	datFile.RomNames = make(map[string]*Game, len(datFile.Games))
	for i := range datFile.Games {
		datFile.RomNames[datFile.Games[i].Name] = &datFile.Games[i]
		datFile.Games[i].RomNames = make(map[string]*Rom, len(datFile.Games[i].Roms))
		datFile.Games[i].Clones = make(map[string]*Game)

		for j := range datFile.Games[i].Roms {
			datFile.Games[i].RomNames[datFile.Games[i].Roms[j].Name] = &datFile.Games[i].Roms[j]
		}
	}

	for i := range datFile.Games {
		if datFile.Games[i].CloneOf != "" {
			datFile.RomNames[datFile.Games[i].CloneOf].Clones[datFile.Games[i].Name] = &datFile.Games[i]
		}
	}

	return &datFile, nil
}
