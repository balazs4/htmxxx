default:
	@test `which templ`     || echo 'go install github.com/a-h/templ/cmd/templ@latest'
	@test `which watchexec` || echo 'https://github.com/watchexec/watchexec?tab=readme-ov-file#install'
watch:
	@WATCH=1 watchexec --print-events --no-meta -c -r --project-origin . --stop-timeout 0 -- make --always-make --silent run

fmt:
	@go fmt

run:
	@templ generate && go run .

htmxxx:
	@templ generate && go build .

htmx.min.js:
	curl https://unpkg.com/htmx.org@1.9.10/dist/htmx.min.js -O
