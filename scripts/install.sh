#!/bin/bash

set -e

PROGRAM_NAME="gas"
GITHUB_USER="style77"
GITHUB_REPO="gas"
INSTALL_DIR="/usr/local/bin"
ARCH=$(uname -m)

# Determine OS type
OS=$(uname -s)

case "$OS" in
    Linux*)
        OS="linux"
        ;;
    Darwin*)
        OS="darwin"
        ;;
    *)
        echo "Unsupported OS: $OS"
        exit 1
        ;;
esac

# Determine architecture
case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    arm64)
        ARCH="arm64"
        ;;
    *)
        echo "Unsupported architecture: $ARCH"
        exit 1
        ;;
esac

echo "Determined OS: $OS and Architecture: $ARCH"

# Fetch latest release from GitHub
echo "Fetching the latest release from GitHub..."
LATEST_RELEASE=$(curl -s "https://api.github.com/repos/$GITHUB_USER/$GITHUB_REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_RELEASE" ]; then
    echo "Error: Could not determine the latest release. Please check the repository name or your internet connection."
    exit 1
fi

echo "Latest release found: $LATEST_RELEASE"

# Construct download URL and filename
FILENAME="${PROGRAM_NAME}-${LATEST_RELEASE}-${OS}-${ARCH}.tar.gz"
DOWNLOAD_URL="https://github.com/$GITHUB_USER/$GITHUB_REPO/releases/download/$LATEST_RELEASE/$FILENAME"

echo "Downloading $FILENAME from $DOWNLOAD_URL..."
curl -L -o "$FILENAME" "$DOWNLOAD_URL"

# Extract the binary
echo "Extracting the binary..."
tar -xzf "$FILENAME"

# Assume that the binary is named after the program
BINARY_NAME="./$PROGRAM_NAME"

if [ ! -f "$BINARY_NAME" ]; then
    echo "Error: Could not find the extracted binary. Please ensure the archive structure is correct."
    exit 1
fi

# Install the binary
echo "Installing $PROGRAM_NAME to $INSTALL_DIR..."

# Make sure the install directory exists
if [ ! -d "$INSTALL_DIR" ]; then
    echo "Creating install directory at $INSTALL_DIR..."
    sudo mkdir -p "$INSTALL_DIR"
fi

# Move binary to the install directory
sudo cp "$BINARY_NAME" "$INSTALL_DIR/"
sudo chmod +x "$INSTALL_DIR/$PROGRAM_NAME"

# Verify if the installation directory is in PATH
if [[ ":$PATH:" == *":$INSTALL_DIR:"* ]]; then
    echo "$PROGRAM_NAME is successfully installed and is already in your PATH."
else
    echo "The installation directory ($INSTALL_DIR) is not in your PATH."
    echo "Consider adding it to your PATH with:"
    echo "export PATH=\$PATH:$INSTALL_DIR"
    echo "You can add this line to your ~/.bashrc or ~/.zshrc file for persistence."
fi

# Clean up
echo "Cleaning up..."
rm "$FILENAME"

echo "$PROGRAM_NAME installation complete. You can now use '$PROGRAM_NAME' in your terminal."
