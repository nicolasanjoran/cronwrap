#!/bin/bash

BINARY_NAME=cronwrap
RELEASE_DIR=release

mkdir -p $RELEASE_DIR

echo "Building for Linux..."
GOOS=linux GOARCH=amd64 go build -o $RELEASE_DIR/${BINARY_NAME}_linux_amd64
GOOS=linux GOARCH=arm go build -o $RELEASE_DIR/${BINARY_NAME}_linux_arm

echo "Building for macOS..."
GOOS=darwin GOARCH=amd64 go build -o $RELEASE_DIR/${BINARY_NAME}_macos_intel
GOOS=darwin GOARCH=arm64 go build -o $RELEASE_DIR/${BINARY_NAME}_macos_arm64

echo "Build complete!"


