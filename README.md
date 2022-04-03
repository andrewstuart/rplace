rplace is a Go client for the r/place canvas and updates.

Its purpose is to be able to more simply diff the state of the canvas with some
desired state as given by an image and x,y location.

```go
var cli rplace.Client
updates, err := cli.Subscribe(context.Background())
if err != nil {
	log.Fatal(err)
}

for upds := range updates {
	for _, upd := range upds {
		fmt.Println(upd.Link(), " => " upd.Color.Name)	
	}
}
```

For a more practical example, see [the cmd/rplace
package](./cmd/rplace/main.go).
