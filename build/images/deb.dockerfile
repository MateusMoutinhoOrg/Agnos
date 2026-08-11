# Build Agnos-Cli and package it as a .deb for Debian/Ubuntu.
# Output: /output/x86-deb.deb
# Installs the binary to /usr/bin/agnos-cli.

FROM golang:latest

WORKDIR /app
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o agnos-cli ./cmd/main

RUN VERSION=$(tr -d '[:space:]' < assets/version.txt) && \
    mkdir -p /pkg/DEBIAN /pkg/usr/bin /output && \
    cp agnos-cli /pkg/usr/bin/agnos-cli && \
    chmod 755 /pkg/usr/bin/agnos-cli && \
    printf 'Package: agnos-cli\nVersion: %s\nArchitecture: amd64\nMaintainer: MateusMoutinhoOrg\nDescription: An OS-independent Go CLI financial tracker template\n' \
        "$VERSION" > /pkg/DEBIAN/control && \
    dpkg-deb --build /pkg /output/x86-deb.deb
