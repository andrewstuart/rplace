package main

import (
	"log"

	"github.com/gopuff/morecontext"
)

func main() {
	ctx := morecontext.ForSignals()

	cli := Client{}
	ch, err := cli.Subscribe(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
