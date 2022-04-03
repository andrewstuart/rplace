package rplace

import (
	"context"
	"fmt"
	"image"
	"image/png"
	"net/http"
)

type basicMessage struct {
	ID      string  `json:"id"`
	Payload payload `json:"payload"`
	Type    string  `json:"type"`
}

func (b basicMessage) getDeltaPng(ctx context.Context) (image.Image, error) {
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

type payload struct {
	Data payloadData `json:"data"`
}

type payloadData struct {
	Subscribe subscribe `json:"subscribe"`
}

type subscribe struct {
	Typename string        `json:"__typename"`
	Data     subscribeData `json:"data"`
	ID       string        `json:"id"`
}

type subscribeData struct {
	Typename          string  `json:"__typename"`
	CurrentTimestamp  float64 `json:"currentTimestamp"`
	Name              string  `json:"name"`
	PreviousTimestamp float64 `json:"previousTimestamp"`
}
