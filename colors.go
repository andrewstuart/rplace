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

var (
	stdPalette, _ = GetPalette(strings.NewReader(colorString))
	canvasColor   = map[[3]uint32]CanvasColor{}
)

func init() {
	for _, c := range stdPalette {
		r, g, b, _ := c.Color.RGBA()
		canvasColor[[3]uint32{r, g, b}] = c
	}
}

func lookupColor(c color.Color) (CanvasColor, error) {
	r, g, b, _ := c.RGBA()
	if cc, ok := canvasColor[[3]uint32{r, g, b}]; ok {
		return cc, nil
	}
	return CanvasColor{}, fmt.Errorf("couldn't find color %s", c)
}

type CanvasColor struct {
	color.Color
	Name string
}

// GetPalette searches for the colors given an html page.
func GetPalette(r io.Reader) (map[string]CanvasColor, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("error parsing doc: %w", err)
	}

	colors := map[string]CanvasColor{}
	doc.Find(".color-container").Each(func(i int, s *goquery.Selection) {
		name := s.Find(".color-name").Text()
		sty, _ := s.Find("button div").Attr("style")

		var hex string
		fmt.Sscanf(sty, "background-color:%s", &hex)
		hex = strings.Split(hex, ";")[0]
		fmt.Printf("hex = %+v\n", hex)

		colors[name] = CanvasColor{
			Color: mustParseHexColor(hex),
			Name:  name,
		}
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

const colorString = `<?xml version="1.0"?>
<div class="container" style="height: 116px;">
  <div class="layout">
    <!--?lit$7247925195$-->
    <div class="palette">
      <!--?lit$7247925195$-->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color" data-color="1">
          <div style="background-color:#BE0039;border:1px solid #BE0039;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->dark red</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="2">
          <div style="background-color:#FF4500;border:1px solid #FF4500;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->red</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="3">
          <div style="background-color:#FFA800;border:1px solid #FFA800;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->orange</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="4">
          <div style="background-color:#FFD635;border:1px solid #FFD635;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->yellow</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="6">
          <div style="background-color:#00A368;border:1px solid #00A368;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->dark green</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="7">
          <div style="background-color:#00CC78;border:1px solid #00CC78;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->green</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="8">
          <div style="background-color:#7EED56;border:1px solid #7EED56;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->light green</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="9">
          <div style="background-color:#00756F;border:1px solid #00756F;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->dark teal</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="10">
          <div style="background-color:#009EAA;border:1px solid #009EAA;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->teal</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="12">
          <div style="background-color:#2450A4;border:1px solid #2450A4;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->dark blue</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="13">
          <div style="background-color:#3690EA;border:1px solid #3690EA;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->blue</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="14">
          <div style="background-color:#51E9F4;border:1px solid #51E9F4;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->light blue</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="15">
          <div style="background-color:#493AC1;border:1px solid #493AC1;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->indigo</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="16">
          <div style="background-color:#6A5CFF;border:1px solid #6A5CFF;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->periwinkle</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="18">
          <div style="background-color:#811E9F;border:1px solid #811E9F;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->dark purple</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="19">
          <div style="background-color:#B44AC0;border:1px solid #B44AC0;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->purple</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="22">
          <div style="background-color:#FF3881;border:1px solid #FF3881;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->pink</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="23">
          <div style="background-color:#FF99AA;border:1px solid #FF99AA;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->light pink</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="24">
          <div style="background-color:#6D482F;border:1px solid #6D482F;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->dark brown</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="25">
          <div style="background-color:#9C6926;border:1px solid #9C6926;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->brown</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="27">
          <div style="background-color:#000000;border:1px solid #000000;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->black</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="29">
          <div style="background-color:#898D90;border:1px solid #898D90;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->gray</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color selected" data-color="30">
          <div style="background-color:#D4D7D9;border:1px solid #D4D7D9;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->light gray</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!---->
      <div class="color-container" style="--num-colors:24;">
        <button class="color   " data-color="31">
          <div style="background-color:#FFFFFF;border:1px solid #E9EBED;"/>
        </button>
        <!--?lit$7247925195$-->
        <mona-lisa-tooltip isopen="" small="" name="">
          <div class="color-name"><!--?lit$7247925195$-->white</div>
        </mona-lisa-tooltip>
      </div>
      <!---->
      <!--?lit$7247925195$-->
    </div>
    <div class="actions">
      <button class="cancel">
        <icon-close/>
      </button>
      <!--?lit$7247925195$-->
      <button class="confirm">
        <icon-checkmark/>
      </button>
    </div>
  </div>
</div>`
