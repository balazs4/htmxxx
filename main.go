package main

import (
	"context"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		index("World").Render(context.Background(), w)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
