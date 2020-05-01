package dat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	f, err := New("test.dat")
	assert.NoError(t, err)

	//file
	assert.Equal(t, "Test Dat", f.Header.Name)
	assert.Equal(t, 2, len(f.Games))
	assert.Equal(t, 2, len(f.RomNames))
	assert.Equal(t, "test.dat", f.Path)
	assert.Equal(t, "romname", f.Games[0].Name)

	//input
	assert.Equal(t, 6, f.Games[0].Input.Buttons)
	assert.Equal(t, "joy8way", f.Games[0].Input.Control)

	//video
	assert.Equal(t, "horizontal", f.Games[0].Video.Orientation)
	assert.Equal(t, "4x3", f.Games[0].Video.AspectRatio())

	//roms
	assert.Equal(t, 2, len(f.Games[0].Roms))
	assert.Equal(t, 2, len(f.Games[0].RomNames))

	//clone
	assert.Empty(t, f.Games[0].CloneOf)
	assert.Equal(t, 1, len(f.Games[0].Clones))
	assert.Equal(t, "romname2", f.Games[1].Name)
	assert.Equal(t, "romname", f.Games[1].RomOf)
	assert.Equal(t, "romname", f.Games[1].CloneOf)
	assert.Equal(t, 2, len(f.Games[1].Roms))

	//romnames map
	assert.Equal(t, "Game (set 2)", f.RomNames["romname2"].Description)
	assert.Equal(t, "Game (set 1)", f.RomNames[f.RomNames["romname2"].CloneOf].Description)
}
