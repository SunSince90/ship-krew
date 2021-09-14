GIT_COMMIT=$(shell git rev-list -1 HEAD)
VERSION=$(shell git describe --tags)

build:
	go build -o bin/users -ldflags "-X main.gitCommit=${GIT_COMMIT} -X main.version=${VERSION}" *.go

run:
	go run main.go