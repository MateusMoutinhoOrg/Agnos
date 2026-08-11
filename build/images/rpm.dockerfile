# Build Agnos-Cli and package it as an .rpm for Fedora/RHEL.
# Output: /output/x86-rpm.rpm
# Installs the binary to /usr/bin/agnos-cli.

# --- Stage 1: compile the Go binary ---
FROM golang:latest AS builder

WORKDIR /app
COPY . .

RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 \
    go build -o /app/agnos-cli ./cmd/main

# --- Stage 2: package the binary into an RPM ---
FROM fedora:latest AS packager

RUN dnf install -y rpm-build && dnf clean all

COPY --from=builder /app/agnos-cli     /root/rpmbuild/SOURCES/agnos-cli
COPY --from=builder /app/assets/version.txt /tmp/version.txt

RUN VERSION=$(tr -d '[:space:]' < /tmp/version.txt) && \
    mkdir -p /root/rpmbuild/{SPECS,BUILD,RPMS,SRPMS,SOURCES} /output && \
    printf '%s\n' \
        "Name:           agnos-cli"                              \
        "Version:        ${VERSION}"                             \
        "Release:        1"                                      \
        "Summary:        An OS-independent Go CLI template"      \
        "License:        Unlicense"                               \
        ""                                                        \
        "%description"                                            \
        "An OS-independent Go CLI financial tracker template."   \
        ""                                                        \
        "%install"                                                \
        "mkdir -p %{buildroot}/usr/bin"                           \
        "cp %{_sourcedir}/agnos-cli %{buildroot}/usr/bin/agnos-cli" \
        ""                                                        \
        "%files"                                                  \
        "/usr/bin/agnos-cli"                                      \
    > /root/rpmbuild/SPECS/agnos-cli.spec && \
    rpmbuild -bb /root/rpmbuild/SPECS/agnos-cli.spec && \
    cp /root/rpmbuild/RPMS/*/*.rpm /output/x86-rpm.rpm
