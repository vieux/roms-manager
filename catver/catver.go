package catver

import (
	"path/filepath"

	"github.com/zieckey/goini"
)

type File struct {
	Categories map[string]string
	Path       string
	ShortPath  string
}

func New(path string) (*File, error) {
	ini := goini.New()
	err := ini.ParseFile(path)
	if err != nil {
		return nil, err
	}

	kvmap, ok := ini.GetKvmap("Category")
	if !ok {
		return nil, nil
	}
	var catverFile File
	catverFile.Path = path
	catverFile.ShortPath = filepath.Join(filepath.Base(filepath.Dir(path)), filepath.Base(path))
	catverFile.Categories = kvmap
	return &catverFile, nil
}
