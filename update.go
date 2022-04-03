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
	return fmt.Sprintf("https://new.reddit.com/r/place/?cx=%d&cy=%d&px=17", u.X, u.Y)
}

// requiresUpdate lets you query whether given a target, and its zero point on
// the canvas, an update should be applied or if it is already within the
// desired parameters.
func (upd Update) getUpdate(tgt image.Image, zero image.Point) (Update, bool) {
	// First, reset the update point from the perspective of the image's zero point.
	inTarget := upd.Sub(zero)
	// Then, if not in bounds, false
	if !inTarget.In(tgt.Bounds()) {
		return Update{}, false
	}

	clr := tgt.At(inTarget.X, inTarget.Y)
	rr, gg, bb, aa := clr.RGBA() // The index inside the target image
	if aa == 0 {
		return Update{}, false
	}

	// Then compare the colors
	r, g, b, _ := upd.Color.RGBA()

	if r == rr && g == gg && b == bb {
		return Update{}, false
	}

	upd.Color = lookupColor(clr)
	return upd, true
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
