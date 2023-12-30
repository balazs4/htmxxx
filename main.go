package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s\n", r.Method, r.URL)
		switch r.Method {
		case "GET":
			index().Render(context.Background(), w)
		case "POST":
			r.ParseForm()
			fmt.Println(r.PostForm)
		default:

		}
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
