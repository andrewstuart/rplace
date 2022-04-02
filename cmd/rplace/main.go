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

	// var o sync.Once
	// colors := map[[3]uint32]color.Color{}
	for up := range ch {
		fmt.Printf("up = %+v\n", up)
		// r, g, b, _ := up.Color.RGBA()
		// k := [3]uint32{r, g, b}
		// colors[k] = up.Color
		// if len(colors) == 24 {
		// 	o.Do(func() {
		// 		for _, v := range colors {
		// 			fmt.Printf("v = %+v\n", v)
		// 		}
		// 	})
		// }
	}
}
