package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

type User struct {
	Name       string
	Email      string
	Validation map[string]string
}

func NewUser(name, email string) *User {
	return &User{
		Name:       name,
		Email:      email,
		Validation: make(map[string]string, 0),
	}
}

func (u *User) IsValid() bool {
	return u.Validation["email"] == "" && u.Validation["name"] == ""
}

func main() {
	var storage = make(map[string]User, 0)
	storage["foo"] = *NewUser("foo", "foo@bar.com")

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
			index(users(&storage, nil)).Render(context.Background(), w)

		case http.MethodPost:
			r.ParseForm()
			user := NewUser(r.PostForm.Get("name"), r.PostForm.Get("email"))

			if htmx := r.Header.Get("hx-request"); htmx == "true" {
				if _, exist := storage[user.Name]; exist == true {
					user.Validation["name"] = "Username is already taken."
				}

				if strings.Contains(user.Email, "@") == false {
					user.Validation["email"] = "Not valid email"
				}

				if user.IsValid() == false {
					users(&storage, user).Render(context.Background(), w)
					return
				}

				storage[user.Name] = *user
				users(&storage, nil).Render(context.Background(), w)
				return
			}

			index(users(&storage, nil)).Render(context.Background(), w)
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
