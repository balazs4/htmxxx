package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	if os.Getenv("WATCH") == "1" {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM)

		http.HandleFunc("/.sigterm", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("content-type", "text/event-stream")
			w.Header().Add("connection", "keep-alive")
			w.Header().Add("cache-control", "no-cache")

			<-sigs
			sigterm := fmt.Sprintf("data: SIGTERM pid %d \n\n", os.Getpid())
			w.Write([]byte(sigterm))
		})
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s\n", r.Method, r.URL)

		switch r.Method {

		case http.MethodGet:
			index().Render(context.Background(), w)

		case http.MethodPost:
			r.ParseForm()
			fmt.Println(r.PostForm)

		}
	})

	http.HandleFunc("/htmx.min.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "htmx.min.js")
	})

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
