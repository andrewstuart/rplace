package main

import (
	"fmt"
	"log"

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

	for upds := range ch {
		for _, up := range upds {
			fmt.Printf("up = %+v\n", up)
		}
	}
}
