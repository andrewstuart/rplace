package rplace

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

// A Client exists to help read and understand changes to the canvas as they
// happen, and determine changes necessary to bring specific desired states
// (images) to the canvas.
type Client struct {
	curr image.Image
}

const knownCanvases = 2

// NeededUpdatesFor takes an image.Image and a location on r/place to place it,
// and returns a channel of updates needed to create and maintain the image at
// that location. It is subscribed to further canvas updates and will continue
// to return necessary changes until the context is closed.
func (c *Client) NeededUpdatesFor(ctx context.Context, img image.Image, at image.Point) (chan Update, error) {
	updch, err := c.Subscribe(ctx)
	if err != nil {
		return nil, fmt.Errorf("error subscribing to updates: %w", err)
	}

	upds := c.getDiff(img, at)
	ch := make(chan Update)
	go func() {
		for _, upd := range upds {
			ch <- upd
		}
	}()
	go func() {
		defer close(ch)

		for upds := range updch {
			for _, upd := range upds {
				if upd.requiresUpdate(c.curr, img, at) {
					select {
					case ch <- upd:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return ch, nil
}

// WithImage returns a copy of the canvas with a preview image overlayed.
func (c Client) WithImage(img image.Image, at image.Point) (image.Image, error) {
	// copy
	buf := &bytes.Buffer{}
	png.Encode(buf, c.curr)
	curr, err := png.Decode(buf)
	if err != nil {
		return nil, fmt.Errorf("error re-decoding for some reason: %w", err)
	}

	draw.Draw(curr.(*image.Paletted), image.Rectangle{Min: at, Max: at.Add(img.Bounds().Max)}, img, image.Point{}, draw.Over)
	return curr, nil
}

// getDiff returns a slice of changes that must be made for the Client's
// current canvas to become the given image, overlayed at the given point.
func (c Client) getDiff(img image.Image, at image.Point) []Update {
	bs := img.Bounds()
	x, y := at.X, at.Y
	var upds []Update

	for xx := 0; xx < bs.Max.X; xx++ {
		for yy := 0; yy < bs.Max.Y; yy++ {
			currColor := c.curr.At(xx+x, yy+y)
			desiredColor := img.At(xx, yy)

			r, g, b, _ := currColor.RGBA()
			rr, gg, bb, _ := desiredColor.RGBA()
			if !(r == rr && g == gg && b == bb) {
				upds = append(upds, Update{
					Point: image.Point{
						X: xx + x,
						Y: yy + y,
					},
					Color: lookupColor(desiredColor),
				})
			}
		}
	}

	return upds
}

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

// getInitial waits for N full frame image updates
func (c *Client) getInitial(ctx context.Context, conn *websocket.Conn, numCanvases int) error {
	for i := 0; i < numCanvases; i++ {
		start.Payload.Variables.Input.Channel.Tag = fmt.Sprint(i)
		err := conn.WriteJSON(start)
		if err != nil {
			return fmt.Errorf("error writing start JSON for canvas #%d: %w", i+1, err)
		}
	}

	i := 0
	c.curr = image.NewPaletted(image.Rect(0, 0, 2000, 1000), stdPalette)
	for i < numCanvases {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		var msg basicMessage
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println(err)
			continue
		}

		if msg.Type != "data" {
			continue
		}

		switch msg.Payload.Data.Subscribe.Data.Typename {
		case "FullFrameMessageData":
			i++
			img, err := msg.getDeltaPng(ctx)
			if err != nil {
				return fmt.Errorf("error getting png: %w", err)
			}

			rect := c.curr.Bounds()
			if strings.Contains(msg.Payload.Data.Subscribe.Data.Name, "-1-") {
				rect.Min.X = 1000
			}
			draw.Draw(c.curr.(*image.Paletted), rect, img, image.Point{}, draw.Over)
		}
	}
	return nil
}

// Subscribe returns a channel of pixel updates from the r/place canvas.
func (c *Client) Subscribe(ctx context.Context) (chan []Update, error) {
	tok, err := GetAnonymousToken(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting anonymous bearer token: %w", err)
	}

	conn, res, err := websocket.DefaultDialer.DialContext(ctx, "wss://gql-realtime-2.reddit.com/query", http.Header{
		"Sec-Websocket-Protocol": []string{"graphql-ws"},
		"Origin":                 []string{"https://hot-potato.reddit.com"},
	})
	if err != nil {
		return nil, fmt.Errorf("error getting websocket conn: %w", err)
	}
	res.Body.Close()

	err = conn.WriteJSON(connectionInitMessage{
		Type: "connection_init",
		Payload: connectionInitMessagePayload{
			Authorization: "Bearer " + tok.AccessToken,
		},
	})

	if err != nil {
		return nil, fmt.Errorf("error authorizing connection: %w", err)
	}

	// Get initial images
	err = c.getInitial(ctx, conn, knownCanvases)
	if err != nil {
		return nil, fmt.Errorf("error getting initial canvases: %w", err)
	}

	ch := make(chan []Update)
	go func() {
		defer conn.Close()
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			var msg basicMessage
			err = conn.ReadJSON(&msg)
			if err != nil {
				log.Println(err)
				continue
			}

			if msg.Type != "data" {
				continue
			}

			img, err := msg.getDeltaPng(ctx)
			if err != nil {
				log.Println(err)
				continue
			}

			rect := c.curr.Bounds()
			if strings.Contains(msg.Payload.Data.Subscribe.Data.Name, "-1-") {
				rect.Min.X = 1000
			}
			draw.Draw(c.curr.(*image.Paletted), rect, img, image.Point{}, draw.Over)

			// Try to send or context closed
			select {
			case <-ctx.Done():
				return
			case ch <- getUpdates(rect.Min, img):
			}
		}

	}()
	return ch, nil
}
