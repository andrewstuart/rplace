package rplace

import (
	"fmt"
	"image/color"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

var (
// DarkRed    = mustParseHexColor("BE0039")
// Red        = mustParseHexColor("FF4500")
// Orange     = mustParseHexColor("FFA800")
// Yellow     = mustParseHexColor("FFD635")
// DarkGreen  = mustParseHexColor("00A368")
// Green      = mustParseHexColor("00CC78")
// LightGreen = mustParseHexColor("7EED56")
// DarkTeal   = mustParseHexColor("00756F")
// Teal       = mustParseHexColor("009EAA")
// DarkBlue   = mustParseHexColor("2450A4")
// Blue       = mustParseHexColor("3690EA")
// LightBlue  = mustParseHexColor("51E9F4")
// Indigo     = mustParseHexColor("493AC1")
// Periwinkle = mustParseHexColor("6A5CFF")
// DarkPurple = mustParseHexColor("811E9F")
// Purple     = mustParseHexColor("B44AC0")
// Pink       = mustParseHexColor("FF3881")
// Pink       = mustParseHexColor("FF3881")
)

func GetPalette(r io.Reader) (map[string]color.Color, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("error parsing doc: %w", err)
	}

	colors := map[string]color.Color{}
	doc.Find(".color-container").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".color-name").Text()
		sty, _ := s.Find("button div").Attr("style")

		var hex string
		fmt.Sscanf(sty, "background-color:%s", &hex)
		hex = strings.Split(hex, ";")[0]
		fmt.Printf("hex = %+v\n", hex)

		colors[name] = mustParseHexColor(hex)
	})
	return colors, nil
}

// https://stackoverflow.com/questions/54197913/parse-hex-string-to-image-color
func mustParseHexColor(s string) (c color.RGBA) {
	var err error
	c.A = 0xff
	switch len(s) {
	case 7:
		_, err = fmt.Sscanf(s, "#%02x%02x%02x", &c.R, &c.G, &c.B)
	case 4:
		_, err = fmt.Sscanf(s, "#%1x%1x%1x", &c.R, &c.G, &c.B)
		// Double the hex digits:
		c.R *= 17
		c.G *= 17
		c.B *= 17
	default:
		err = fmt.Errorf("invalid length, must be 7 or 4")

	}
	if err != nil {
		panic(err)
	}
	return c
}
