FROM golang:1.17 as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go main.go
COPY create_user.go create_user.go
COPY delete_user.go delete_user.go
COPY get_user.go get_user.go
COPY list_users.go list_users.go
COPY update_user.go update_user.go
COPY pkg/ pkg/

RUN export GIT_COMMIT=$(shell git rev-list -1 HEAD)

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GO111MODULE=on \
    go build -a -o users-api-server -ldflags "-X main.gitCommit=${GIT_COMMIT}" *.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/users-api-server .
USER nonroot:nonroot

ENTRYPOINT ["/users-api-server"]
