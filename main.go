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
	http.Handle("/htmx.min.js", http.FileServer(http.FS(htmx)))

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	watch := os.Getenv("WATCH") == "1"
	if watch == true {
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

	var storage = make(map[string]structs.User, 0)
	storage["foo"] = *structs.NewUser("foo", "foo@bar.com")

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s\n", r.Method, r.URL)
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		hx_request := r.Header.Get("hx-request") == "true"

		r.ParseForm()

		user := structs.NewUser(r.PostForm.Get("name"), r.PostForm.Get("email"))
		user.Validate(&storage)

		if user.IsValid() == true {
			storage[user.Name] = *user

			if hx_request == true {
				if err := html.Page.ExecuteTemplate(w, "main", html.MainProps{nil, &storage}); err != nil {
					fmt.Println(err)
				}
				return
			}

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
			return
		}

		if hx_request == true {
			if err := html.Page.ExecuteTemplate(w, "main", html.MainProps{user, &storage}); err != nil {
				fmt.Println(err)
			}
			return
		}

		if err := html.Page.ExecuteTemplate(w, "page", html.PageProps{watch, html.MainProps{user, &storage}}); err != nil {
			fmt.Println(err)
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s\n", r.Method, r.URL)

		if err := html.Page.ExecuteTemplate(w, "page", html.PageProps{watch, html.MainProps{nil, &storage}}); err != nil {
			fmt.Println(err)
		}
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
