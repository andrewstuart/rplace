package main

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"net/http"
)

type BasicMessage struct {
	ID      string  `json:"id"`
	Payload Payload `json:"payload"`
	Type    string  `json:"type"`
}

func (b BasicMessage) getDeltaPng(ctx context.Context) (image.Image, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, b.Payload.Data.Subscribe.Data.Name, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating http request for image: %w", err)
	}

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request for image: %w", err)
	}

	defer res.Body.Close()
	return png.Decode(res.Body)
}

type Payload struct {
	Data PayloadData `json:"data"`
}

type PayloadData struct {
	Subscribe Subscribe `json:"subscribe"`
}

type Subscribe struct {
	Typename string        `json:"__typename"`
	Data     SubscribeData `json:"data"`
	ID       string        `json:"id"`
}

type SubscribeData struct {
	Typename          string  `json:"__typename"`
	CurrentTimestamp  float64 `json:"currentTimestamp"`
	Name              string  `json:"name"`
	PreviousTimestamp float64 `json:"previousTimestamp"`
}
