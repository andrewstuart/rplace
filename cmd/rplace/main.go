package main

import (
	"fmt"
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

	ch, err := cli.NeededUpdatesFor(ctx, img, 906, 646)
	if err != nil {
		log.Fatal(err)
	}

	for up := range ch {
		fmt.Printf("Visit %s and select %s\n", up.Link(), up.Color.Name)
	}
}
