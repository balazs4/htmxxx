package main

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

//go:embed htmx.min.js
var htmx embed.FS

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
	return u.Validation["Email"] == "" && u.Validation["Name"] == ""
}

type std_index_props struct {
	Body  template.HTML
	Watch bool
}

var std_index = template.Must(template.New("index").Parse(`
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8"/>
    <title>go + templ + htmx</title>
    <script src="/htmx.min.js"></script>
    <style>
    @media (prefers-color-scheme: light) {
      * { background: #ffffff; color: #000000; }
    }
    @media (prefers-color-scheme: dark) {
      * { background: #000000; color: #ffffff; }
      input { background: #333; color: #eee; }
    }
    </style>
  </head>
  <body>
  {{ .Body }}
  </body>
</html>
{{ if .Watch }}
<script>
  new EventSource("/.sigterm").onmessage = function(){ setTimeout(() => { location.reload(); }, 750); }
</script>
{{ end }}
`))

type std_users_props struct {
	User  *User
	Users *map[string]User
}

var std_users = template.Must(template.New("users").Parse(`
	<div id="users">
		<form
			hx-trigger="submit"
			hx-post="/register"
			hx-swap="outerHTML"
			hx-target="#users"
			hx-on::after-request=" if(event.detail.successful) this.reset()"
		>
			<label for="name">Username:</label>
      <input type="text" name="name" {{ if .User }}value="{{ .User.Name }}"{{ end }}/>
      {{ if and .User .User.Validation.Name }} <span style="color: red;">{{ .User.Validation.Name}}</span> {{ end }}

			<label for="email">Email:</label>
      <input type="email" name="email" {{ if .User }}value="{{ .User.Email }}"{{ end }}/>
      {{ if and .User .User.Validation.Email }} <span style="color: red;">{{ .User.Validation.Email }}</span> {{ end }}

			<input type="submit" value="BOB!"/>
		</form>
		<ul>
    {{ range .Users }}
    <li><p>{{ .Name }} <a href="mailto:{{ .Email }}">{{ .Email }}</a></p></li>
    {{ end }}
		</ul>
	</div>
`))

func main() {
	watch := os.Getenv("WATCH") == "1"

	var storage = make(map[string]User, 0)
	storage["foo"] = *NewUser("foo", "foo@bar.com")

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

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s\n", r.Method, r.URL)

		switch r.Method {

		case http.MethodGet:
			var buff bytes.Buffer
			std_users.Execute(&buff, std_users_props{User: nil, Users: &storage})
			std_index.Execute(w,
				std_index_props{
					Body:  template.HTML(buff.String()),
					Watch: watch,
				})

		case http.MethodPost:
			r.ParseForm()
			user := NewUser(r.PostForm.Get("name"), r.PostForm.Get("email"))

			if htmx := r.Header.Get("hx-request"); htmx == "true" {
				if _, exist := storage[user.Name]; exist == true {
					user.Validation["Name"] = "Username is already taken."
				}

				if strings.Contains(user.Email, "@") == false {
					user.Validation["Email"] = "Not valid email"
				}

				if user.IsValid() == false {
					std_users.Execute(w, std_users_props{user, &storage})
					return
				}

				storage[user.Name] = *user
				std_users.Execute(w, std_users_props{nil, &storage})
				return
			}

			http.Redirect(w, r, "/", http.StatusMovedPermanently)
		}
	})

	http.Handle("/htmx.min.js", http.FileServer(http.FS(htmx)))

	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})

	log.Fatal(http.ListenAndServe(":3000", nil))
}
