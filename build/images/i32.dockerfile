# Build Agnos-Cli for Windows i386.
# Output: /output/windows-i32.exe

FROM golang:latest

WORKDIR /app
COPY . .

RUN mkdir -p /output && \
    GOOS=windows GOARCH=386 CGO_ENABLED=0 \
    go build -o /output/windows-i32.exe ./cmd/main
