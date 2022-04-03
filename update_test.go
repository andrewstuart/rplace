package rplace

import (
	"image"
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRequiresUpdate(t *testing.T) {
	asrt := assert.New(t)

	sub := image.NewPaletted(image.Rect(0, 0, 100, 100), stdPalette)
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	sub.Set(0, 0, white)

	u := Update{
		Point: image.Point{0, 0},
		Color: CanvasColor{Color: color.RGBA{}},
	}

	up, ok := u.getUpdate(sub, image.Point{})
	asrt.True(ok)
	asrt.Equal(up.Color.Color, white)

	// Bounds checking
	u.Point.X = 100
	_, ok = u.getUpdate(sub, image.Point{})
	asrt.False(ok)
	_, ok = u.getUpdate(sub, image.Point{X: 1})
	asrt.True(ok)

	// Now test we don't create Updates for something with transparency in the target
	sub.Set(1, 1, color.RGBA{R: 0, G: 0, B: 255, A: 0})
	u.Point = image.Point{X: 1, Y: 1}

	sub.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 0})
	_, ok = u.getUpdate(sub, image.Point{})
	asrt.False(ok)
}
