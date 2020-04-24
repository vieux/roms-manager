package dat

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"
)

type File struct {
	Header Header `xml:"header"`
	Games  []Game `xml:"game"`

	Map       map[string]*Game `xml:"-"`
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
	Rom    []Rom `xml:"rom"`
	Driver struct {
		Status string `xml:"status,attr"`
	} `xml:"driver"`
	Map map[string]*Rom `xml:"-"`
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
	defer f.Close()
	if err != nil {
		return nil, err
	}

	var datFile File
	if err := xml.NewDecoder(f).Decode(&datFile); err != nil {
		return nil, err
	}

	datFile.Path = path
	datFile.ShortPath = filepath.Join(filepath.Base(filepath.Dir(path)), filepath.Base(path))

	datFile.Map = make(map[string]*Game, len(datFile.Games))
	for i := range datFile.Games {
		datFile.Map[datFile.Games[i].Name] = &datFile.Games[i]
		datFile.Games[i].Map = make(map[string]*Rom, len(datFile.Games[i].Rom))
		for j := range datFile.Games[i].Rom {
			ext := filepath.Ext(datFile.Games[i].Rom[j].Name)
			datFile.Games[i].Map[strings.TrimSuffix(datFile.Games[i].Rom[j].Name, ext)] = &datFile.Games[i].Rom[j]
		}
	}

	return &datFile, nil
}
