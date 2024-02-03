package html

import (
	"html/template"

	"github.com/balazs4/htmxxx/structs"
)

type MainProps struct {
	User  *structs.User
	Users *map[string]structs.User
}

type PageProps struct {
	Watch bool
	Main  MainProps
}

var Page = template.Must(template.New("page").Parse(
	`
{{ if .Watch }}
<script> new EventSource("/.sigterm").onmessage = function(){ setTimeout(() => { location.reload(); }, 750); }</script>
{{ end }}
<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="UTF-8"/>
    <title>users</title>
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
  <header>{{ block "header" . }}{{end}}</header>

  <main>{{ block "main" .Main }}

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

			<input type="submit" value="Register"/>
		</form>
		<ul>
    {{ range .Users }}
    <li><p>{{ .Name }} <a href="mailto:{{ .Email }}">{{ .Email }}</a></p></li>
    {{ end }}
		</ul>
	</div>
  {{end}}</main>

  <footer>{{ block "footer" . }}{{end}}</footer>
  </body>
</html>
`))
