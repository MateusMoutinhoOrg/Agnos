# Build Agnos-Cli for macOS Intel (amd64).
# Output: /output/mac-amd64.bin

FROM golang:latest

WORKDIR /app
COPY . .

RUN mkdir -p /output && \
    GOOS=darwin GOARCH=amd64 CGO_ENABLED=0 \
    go build -o /output/mac-amd64.bin ./cmd/main
