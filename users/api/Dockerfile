# Build the binary.
FROM golang:1.17 as builder

WORKDIR /workspace

# Copy the Go Modules manifests.
COPY go.mod go.mod
COPY go.sum go.sum

# Cache deps before building and copying source so that we don't need to
# re-download as much and so that source changes don't invalidate our
# downloaded layer.
RUN go mod download

# Copy the go source.
COPY main.go main.go
COPY internal/ internal/
COPY pkg/ pkg/

# Build, based on the architecture we want this to run.
ENV ARCH=amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=${ARCH} GO111MODULE=on \
go build -a -o users-api main.go

# Use distroless as minimal base image to package the binary.
# Refer to https://github.com/GoogleContainerTools/distroless for more details.
FROM gcr.io/distroless/static:nonroot
WORKDIR /
COPY --from=builder /workspace/users-api .
USER nonroot:nonroot

LABEL app=users
LABEL module=api

EXPOSE 8080 8081
ENTRYPOINT ["/users-api"]