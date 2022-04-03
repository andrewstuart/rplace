package rplace

import (
	"fmt"
	"image"
)

// An Update is a desired change that must be made to a canvas.
type Update struct {
	image.Point
	Color CanvasColor
}

// Link encodes the browser link that would take you to the correct location on
// r/place's canvas.
func (u Update) Link() string {
	return fmt.Sprintf("https://www.reddit.com/r/place/?cx=%d&cy=%d&px=17", u.X, u.Y)
}

// requiresUpdate lets you query whether given a canvas, target, and point, an
// update should be applied or if it is already within the desired parameters.
func (upd Update) requiresUpdate(canvas, tgt image.Image, pt image.Point) bool {
	bs := tgt.Bounds()
	if !pt.In(bs) {
		return false
	}

	r, g, b, _ := upd.Color.RGBA()
	rr, gg, bb, _ := stdPalette.Convert(tgt.At(upd.X-pt.X, upd.Y-pt.Y)).RGBA() // The index inside the target image

	return !(r == rr && g == gg && b == bb)
}

// GetUpdates returns the list of updated pixels from an image, given that
// image's known zero point on the canvas. The zero point was added after day 2
// when the reddit r/place canvas doubled in size.
func getUpdates(zero image.Point, img image.Image) []Update {
	var upd []Update
	bs := img.Bounds()
	for i := 0; i < bs.Max.X; i++ {
		for j := 0; j < bs.Max.Y; j++ {
			clr := img.At(i, j)
			_, _, _, a := clr.RGBA()
			if a > 0 {
				upd = append(upd, Update{
					Point: zero.Add(image.Point{X: i, Y: j}),
					Color: lookupColor(clr),
				})
			}
		}
	}
	return upd
}
