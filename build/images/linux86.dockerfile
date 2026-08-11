# Build Agnos-Cli as a static Linux amd64 binary.
# Output: /output/linux86.out

FROM golang:latest

WORKDIR /app
COPY . .

RUN mkdir -p /output && \
    GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o /output/linux86.out ./cmd/main
