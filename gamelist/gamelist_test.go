package gamelist

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	f, err := New("test.xml")
	assert.NoError(t, err)

	//file
	assert.Equal(t, 2, len(f.Games))
	assert.Equal(t, 2, len(f.RomNames))
	assert.Equal(t, "test.xml", f.Path)
	assert.Equal(t, "romname", f.Games[0].RomName)
	assert.True(t, f.Games[0].Hidden)
	assert.False(t, f.Games[1].Hidden)

	//romnames map
	assert.Equal(t, "Game (set 1)", f.RomNames["romname"].Name)
	assert.Equal(t, "Game (set 2)", f.RomNames["romname2"].Name)
}
