package html

import (
	"html/template"
	"io"
)

type IndexProps struct {
	Body  template.HTML
	Watch bool
}

var index = template.Must(template.New("index").Parse(
	`{{ if .Watch }}
<script> new EventSource("/.sigterm").onmessage = function(){ setTimeout(() => { location.reload(); }, 750); }</script>
{{ end }}
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
`))

func Index(w io.Writer, p IndexProps) error {
	return index.Execute(w, p)
}
