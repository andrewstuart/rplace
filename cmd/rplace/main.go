package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"

	"github.com/andrewstuart/rplace"
	"github.com/gopuff/morecontext"
)

func main() {
	ctx := morecontext.ForSignals()

	cli := rplace.Client{}

	f, err := os.OpenFile("gopher.png", os.O_RDONLY, 0400)
	if err != nil {
		log.Fatal(err)
	}
	img, err := png.Decode(f)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	ch, err := cli.NeededUpdatesFor(ctx, img, 1434, 664)
	if err != nil {
		log.Fatal(err)
	}

	out, _ := cli.WithImage(img, image.Point{X: 1446, Y: 648})
	f, err = os.OpenFile("test.png", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatal(err)
	}
	png.Encode(f, out)
	if err != nil {
		log.Fatal(err)
	}
	f.Close()

	for up := range ch {
		fmt.Printf("Visit %s and select %s\n", up.Link(), up.Color.Name)
	}
}
