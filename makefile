dev:
	@watchexec --print-events -c -k -r --project-origin . -- make --silent --always-make run

fmt:
	@go fmt

run:
	@go run main.go

build:
	@go build .
