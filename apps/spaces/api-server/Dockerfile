FROM golang:1.17 as builder

WORKDIR /workspace

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod download

COPY main.go main.go
COPY pkg/ pkg/

RUN export GIT_COMMIT=$(shell git rev-list -1 HEAD)

RUN CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm \
    GOARM=7 \
    GO111MODULE=on \
    go build -a -o spaces-api-server -ldflags "-X main.gitCommit=${GIT_COMMIT}" *.go

FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/spaces-api-server .
USER nonroot:nonroot

ENTRYPOINT ["/spaces-api-server"]