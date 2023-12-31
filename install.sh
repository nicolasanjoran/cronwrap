#!/bin/sh

# Set the host URL
HOST_URL="https://raw.githubusercontent.com/nicolasanjoran/cronwrap/main/release"

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
response_code=$(curl -sL -w "%{http_code}" -o /tmp/$FILE_NAME $DOWNLOAD_URL)

if [ "$response_code" -eq 200 ]; then
    echo "Downloaded binary to /tmp/$FILE_NAME"
else
    echo "Received HTTP code: $response_code, cannot download binary"
    exit 1
fi

chmod +x /tmp/$FILE_NAME

# Check if sudo is needed
MOVE_COMMAND="mv /tmp/$FILE_NAME /usr/local/bin/cronwrap"
if [ "$(id -u)" != "0" ]; then
    if command -v sudo > /dev/null 2>&1; then
        MOVE_COMMAND="sudo $MOVE_COMMAND"
    else
        echo "Error: sudo command is not available, but is required to install the binary. Please run the script as root or install sudo."
        exit 1
    fi
fi

$MOVE_COMMAND

if [ $? -ne 0 ]; then
    echo "Error moving ${FILE_NAME} to /usr/local/bin. Do you have the right permissions?"
    exit 1
fi

echo "${FILE_NAME} is now installed and can be run with 'cronwrap'"
