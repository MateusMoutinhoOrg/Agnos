# Build Agnos-Cli for macOS Apple Silicon (arm64).
# Output: /output/mac-arm64.bin

FROM golang:latest

WORKDIR /app
COPY . .

RUN mkdir -p /output && \
    GOOS=darwin GOARCH=arm64 CGO_ENABLED=0 \
    go build -o /output/mac-arm64.bin ./cmd/main
