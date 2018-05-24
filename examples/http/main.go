package main

import (
	"log"
	"net/http"

	"github.com/deadcheat/goblet/examples/http/assetsbin"
)

func main() {
	http.Handle("/static/", http.FileServer(assetsbin.Assets.WithPrefix("/static/")))
	log.Println("start server localhost:3000")
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatal(err)
	}
}
