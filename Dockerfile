FROM golang:1.22 as builder

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . /workspace

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build  -a -o app -v cmd/server/main.go

FROM alpine:latest

COPY --from=builder /workspace/app /usr/bin/

ENTRYPOINT ["/usr/bin/app"]