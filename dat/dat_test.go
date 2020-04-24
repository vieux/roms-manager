package dat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	f, err := New("test.dat")
	assert.NoError(t, err)
	assert.Equal(t, "Test Dat", f.Header.Name)
	assert.Equal(t, 2, len(f.Games))
	assert.Equal(t, "2020bb", f.Games[0].Name)
	assert.Equal(t, "neogeo", f.Games[0].RomOf)
	assert.Equal(t, "horizontal", f.Games[0].Video.Orientation)

	assert.Equal(t, 46, len(f.Games[0].Rom))
	assert.Empty(t, f.Games[0].CloneOf)
	assert.Equal(t, "2020bba", f.Games[1].Name)
	assert.Equal(t, "2020bb", f.Games[1].RomOf)
	assert.Equal(t, "2020bb", f.Games[1].CloneOf)
	assert.Equal(t, 46, len(f.Games[1].Rom))

	assert.Equal(t, 2, len(f.Map))

	assert.Equal(t, "2020 Super Baseball (set 1)", f.Map[f.Map["2020bba"].CloneOf].Description)
}
