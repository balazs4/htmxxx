default:
	@test `which templ`     || echo 'go install github.com/a-h/templ/cmd/templ@latest'
	@test `which watchexec` || echo 'https://github.com/watchexec/watchexec?tab=readme-ov-file#install'
dev:
	@watchexec --print-events -c -r --project-origin . -- make --silent --always-make run

fmt:
	@go fmt

run:
	@templ generate
	@go run .

build:
	@templ generate
	@go build

