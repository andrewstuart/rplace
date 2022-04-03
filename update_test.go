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

	sub.Set(0, 0, color.RGBA{R: 255, G: 255, B: 255, A: 255})

	u := Update{
		Point: image.Point{0, 0},
		Color: CanvasColor{Color: darkRed},
	}

	asrt.True(u.requiresUpdate(sub, image.Point{}))

	// Bounds checking
	u.Point.X = 100
	asrt.False(u.requiresUpdate(sub, image.Point{}))
	asrt.True(u.requiresUpdate(sub, image.Point{X: 1}))

	// Now test we don't create Updates for something with transparency in the target
	sub.Set(1, 1, color.RGBA{R: 0, G: 0, B: 255, A: 0})
	u.Point = image.Point{X: 1, Y: 1}

	sub.Set(1, 1, color.RGBA{R: 255, G: 255, B: 255, A: 0})
	asrt.False(u.requiresUpdate(sub, image.Point{}))

}
