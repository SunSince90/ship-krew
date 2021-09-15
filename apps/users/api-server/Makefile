IMG ?= <repository>
GIT_COMMIT=$(shell git rev-list -1 HEAD)
# VERSION=$(shell git describe --tags)

build:
	go build -o bin/users -ldflags "-X main.gitCommit=${GIT_COMMIT}" *.go

run:
	go run main.go

docker-build:
	docker build . -t ${IMG}

# Push the docker image
docker-push:
	docker push ${IMG}