package rplace

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// A client exists to help read and understand changes to the canvas as they
// happen.
type Client struct {
	curr image.Image
}

// NeededUpdatesFor takes an image.Image and a location on r/place to place it,
// and returns a channel of updates needed to create and maintain the image at
// that location. It is subscribed to further canvas updates and will continue
// to return necessary changes until the context is closed.
func (c Client) NeededUpdatesFor(ctx context.Context, img image.Image, x, y int) (chan Update, error) {
	updch, err := c.Subscribe(ctx)
	if err != nil {
		return nil, fmt.Errorf("error subscribing to updates: %w", err)
	}

	upds := c.getDiff(img, x, y)
	ch := make(chan Update)
	go func() {
		for _, upd := range upds {
			ch <- upd
		}
	}()
	go func() {
		defer close(ch)

		bs := img.Bounds()
		for upds := range updch {
			for _, upd := range upds {
				if upd.X > x && upd.X <= bs.Max.X+x && upd.Y > y && upd.Y <= y+bs.Max.Y {
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

func (c Client) getDiff(img image.Image, x, y int) []Update {
	bs := img.Bounds()
	var upds []Update

	for xx := 0; xx < bs.Max.X; xx++ {
		for yy := 0; yy < bs.Max.Y; yy++ {
			currColor := c.curr.At(xx+x, yy+y)
			desiredColor := img.At(xx, yy)

			r, g, b, _ := currColor.RGBA()
			rr, gg, bb, _ := desiredColor.RGBA()
			if !(r == rr && g == gg && b == bb) {
				upds = append(upds, Update{
					X:     xx + x,
					Y:     yy + y,
					Color: desiredColor,
				})
			}
		}
	}

	return upds
}

type Update struct {
	X, Y  int
	Color color.Color
}

func (upd Update) requiresUpdate(canvas, tgt image.Image, x, y int) bool {
	bs := tgt.Bounds()
	if !(upd.X > x && upd.X <= bs.Max.X+x && upd.Y > y && upd.Y <= y+bs.Max.Y) {
		return false
	}

	r, g, b, _ := upd.Color.RGBA()
	rr, gg, bb, _ := tgt.At(upd.X-x, upd.Y-y).RGBA() // The index inside the target image

	return !(r == rr && g == gg && b == bb)
}

func getUpdates(img image.Image) []Update {
	var upd []Update
	bs := img.Bounds()
	for i := 0; i < bs.Max.X; i++ {
		for j := 0; j < bs.Max.Y; j++ {
			color := img.At(i, j)
			_, _, _, a := color.RGBA()
			if a > 0 {
				upd = append(upd, Update{
					X:     i,
					Y:     j,
					Color: color,
				})
			}
		}
	}
	return upd
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

	err = conn.WriteJSON(start)
	if err != nil {
		return nil, fmt.Errorf("error writing start JSON: %w", err)
	}

	for {
		var msg basicMessage
		err = conn.ReadJSON(&msg)
		if err != nil {
			return nil, fmt.Errorf("error: %w", err)
		}

		if msg.Type != "data" {
			continue
		}

		c.curr, err = msg.getDeltaPng(ctx)
		if err != nil {
			continue
		}
		break
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

			draw.Draw(c.curr.(*image.Paletted), c.curr.Bounds(), img, image.Point{}, draw.Over)

			// Try to send or context closed
			select {
			case <-ctx.Done():
				return
			case ch <- getUpdates(img):
			}
		}

	}()
	return ch, nil
}
