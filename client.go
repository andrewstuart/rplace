package main

import (
	"context"
	"fmt"
	"image"
	"image/color"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	img image.Image
}

type Update struct {
	X, Y  int
	Color color.Color
}

func (c Client) Subscribe(ctx context.Context) (chan []Update, error) {
	tok, err := getToken(ctx)
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
	defer conn.Close()

	err = conn.WriteJSON(ConnectionInitMessage{
		Type: "connection_init",
		Payload: ConnectionInitMessagePayload{
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

	ch := make(chan BasicMessage)
	go func() {
		defer close(ch)

		for {
			select {
			case <-ctx.Done():
				return
			default:
			}

			var msg BasicMessage
			err = conn.ReadJSON(&msg)
			if err != nil {
				return
			}

			var upd []Update

			img, err := msg.getDeltaPng(ctx)
			if err != nil {
				log.Println(err)
				continue
			}

			// Try to send or context closed
			select {
			case <-ctx.Done():
				return
			case ch <- msg:
			}
		}

	}()
	return ch, nil
}
