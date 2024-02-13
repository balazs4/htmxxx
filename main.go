package main

import (
	"embed"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/balazs4/htmxxx/html"
	"github.com/balazs4/htmxxx/structs"
)

//go:embed htmx.min.js
var htmx embed.FS

func main() {
	watch := os.Getenv("WATCH") == "1"
	if watch == true {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGTERM)

		http.HandleFunc("GET /.sigterm", func(w http.ResponseWriter, r *http.Request) {
			w.Header().Add("content-type", "text/event-stream")
			w.Header().Add("connection", "keep-alive")
			w.Header().Add("cache-control", "no-cache")

			<-sigs

			sigterm := fmt.Sprintf("data: SIGTERM pid %d \n\n", os.Getpid())
			w.Write([]byte(sigterm))
		})
	}

	var storage = make(map[string]structs.User, 0)
	storage["foo"] = *structs.NewUser("foo", "foo@bar.com")

	http.Handle("GET /htmx.min.js", http.FileServer(http.FS(htmx)))
	http.HandleFunc("GET /favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	http.HandleFunc("POST /register", func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		user := structs.NewUser(r.PostForm.Get("name"), r.PostForm.Get("email"))
		user.Validate(&storage)

		if user.IsValid() == false {
			props := html.MainProps{user, &storage}
			if err := html.Page.ExecuteTemplate(w, "main", props); err != nil {
				fmt.Println(err)
			}
      return
		}

		storage[user.Name] = *user
		props := html.MainProps{nil, &storage}
		if err := html.Page.ExecuteTemplate(w, "main", props); err != nil {
			fmt.Println(err)
		}
	})

	http.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		props := html.PageProps{watch, html.MainProps{nil, &storage}}
		if err := html.Page.ExecuteTemplate(w, "page", props); err != nil {
			fmt.Println(err)
		}
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
