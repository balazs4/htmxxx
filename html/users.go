package html

import (
	"html/template"
	"io"

	"github.com/balazs4/htmxxx/types"
)

type UsersProps struct {
	User  *types.User
	Users *map[string]types.User
}

var users = template.Must(template.New("users").Parse(`
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
`))

func Users(w io.Writer, p UsersProps) error {
	return users.Execute(w, p)
}
