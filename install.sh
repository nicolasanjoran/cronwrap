#!/bin/sh

# Set the host URL
HOST_URL="https://github.com/nicolasanjoran/cronwrap/raw/main/release"

# Detect OS and Architecture
OS="$(uname | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

# Map architecture to expected names
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    arm64)
        ARCH="arm64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    arm*)
        ARCH="arm"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

# Build download URL
if [ "$OS" = "darwin" ]; then
    OS="macos"
fi

FILE_NAME="cronwrap_${OS}_${ARCH}"
DOWNLOAD_URL="${HOST_URL}/${FILE_NAME}"
echo $DOWNLOAD_URL

# Fetch the binary
curl -Lo /tmp/$FILE_NAME $DOWNLOAD_URL

if [ $? -ne 0 ]; then
    echo "Error downloading ${FILE_NAME}. Please ensure it exists on the server."
    exit 1
fi

chmod +x /tmp/$FILE_NAME

# Move to /usr/local/bin or another location in PATH
sudo mv /tmp/$FILE_NAME /usr/local/bin/cronwrap

if [ $? -ne 0 ]; then
    echo "Error moving ${FILE_NAME} to /usr/local/bin. Do you have the right permissions?"
    exit 1
fi

echo "${FILE_NAME} is now installed and can be run with 'cronwrap'"