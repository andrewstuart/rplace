package main

import (
	"fmt"
	"image/color"
	"log"
	"sync"

	"github.com/andrewstuart/rplace"
	"github.com/gopuff/morecontext"
)

func main() {
	ctx := morecontext.ForSignals()

	cli := rplace.Client{}
	ch, err := cli.Subscribe(ctx)
	if err != nil {
		log.Fatal(err)
	}

	var o sync.Once
	colors := map[[3]uint32]color.Color{}
	for upds := range ch {
		for _, up := range upds {
			r, g, b, _ := up.Color.RGBA()
			k := [3]uint32{r, g, b}
			colors[k] = up.Color
			if len(colors) == 24 {
				o.Do(func() {
					for _, v := range colors {
						fmt.Printf("v = %+v\n", v)
					}
				})
			}
		}
	}
}
