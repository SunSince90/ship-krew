GIT_COMMIT=$(shell git rev-list -1 HEAD)
# VERSION=$(shell git describe --tags)

build:
	go build -o bin/spaces-api-server -ldflags "-X main.gitCommit=${GIT_COMMIT}" *.go

run:
	go run main.go