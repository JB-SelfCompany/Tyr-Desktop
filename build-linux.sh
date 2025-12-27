#!/bin/bash
set -e

# Ensure /usr/local/go/bin is in PATH
export PATH=$PATH:/usr/local/go/bin

echo "========================================"
echo "Tyr Desktop - Linux Build Script"
echo "========================================"
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "ERROR: Go is not installed or not in PATH!"
    echo "Please install Go from https://golang.org/dl/"
    echo "Minimum required version: 1.21"
    echo ""
    echo "For Ubuntu/Debian:"
    echo "  wget https://go.dev/dl/go1.23.0.linux-amd64.tar.gz"
    echo "  sudo tar -C /usr/local -xzf go1.23.0.linux-amd64.tar.gz"
    echo "  export PATH=\$PATH:/usr/local/go/bin"
    exit 1
fi

# Display Go version
echo "Checking Go installation..."
go version
echo ""

# Check if Node.js/npm is installed
if ! command -v node &> /dev/null; then
    echo "ERROR: Node.js is not installed or not in PATH!"
    echo "Please install Node.js from https://nodejs.org/"
    echo ""
    echo "For Ubuntu/Debian:"
    echo "  curl -fsSL https://deb.nodesource.com/setup_lts.x | sudo -E bash -"
    echo "  sudo apt-get install -y nodejs"
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo "ERROR: npm is not installed!"
    echo "Please install npm (usually comes with Node.js)"
    exit 1
fi

# Display Node.js and npm versions
echo "Checking Node.js installation..."
node --version
npm --version
echo ""

# Install/update Wails CLI if not installed
echo "Checking Wails CLI installation..."
if ! command -v wails &> /dev/null; then
    echo "Wails CLI is not installed. Installing now..."
    echo "This may take a few minutes..."
    go install github.com/wailsapp/wails/v2/cmd/wails@latest

    # Add GOPATH/bin to PATH if not already there
    if [[ ":$PATH:" != *":$HOME/go/bin:"* ]]; then
        export PATH=$PATH:$HOME/go/bin
        echo "Added $HOME/go/bin to PATH for this session"
        echo "To make this permanent, add this line to your ~/.bashrc or ~/.zshrc:"
        echo "  export PATH=\$PATH:\$HOME/go/bin"
        echo ""
    fi

    # Verify installation
    if ! command -v wails &> /dev/null; then
        echo "ERROR: Failed to install Wails CLI"
        echo "Please check your Go installation and try again"
        echo "You may need to add \$HOME/go/bin to your PATH"
        exit 1
    fi

    echo "Wails CLI installed successfully!"
    echo ""
else
    echo "Wails CLI is already installed"
    wails version
    echo ""
fi

# Download Go module dependencies
echo "Downloading Go module dependencies..."
go mod download
if [ $? -ne 0 ]; then
    echo "ERROR: Failed to download Go dependencies"
    exit 1
fi
echo "Go dependencies downloaded successfully!"
echo ""

# Extract version from wails.json
echo "Reading version from wails.json..."
VERSION=$(grep -oP '"version":\s*"\K[^"]+' wails.json)
echo "Building version: $VERSION"
echo ""

# Step 1: Preparing build directory...
echo "Step 1: Preparing build directory..."
# Clean build directories
rm -rf build/bin
rm -rf build/linux
rm -rf frontend/dist
echo "Cleaned previous builds."

# Create build directory structure
mkdir -p build/linux

# Copy build assets from assets/build
echo "Copying build assets from assets/build..."
if [ -f "assets/build/appicon.png" ]; then
    cp assets/build/appicon.png build/appicon.png
    echo "  - Copied appicon.png"
else
    echo "  WARNING: assets/build/appicon.png not found!"
fi
if [ -f "assets/build/linux/icon.png" ]; then
    cp assets/build/linux/icon.png build/linux/icon.png
    echo "  - Copied icon.png"
else
    echo "  WARNING: assets/build/linux/icon.png not found!"
fi
echo "Done."
echo ""

# Step 2: Install frontend dependencies
echo "Step 2: Installing frontend dependencies..."
cd frontend
npm install
cd ..
echo "Done."
echo ""

# Step 3: Building application...
echo "Step 3: Building application..."
echo "This may take several minutes..."
wails build -clean -platform linux/amd64 -tags webkit2_41 -ldflags "-X main.version=$VERSION"
if [ $? -ne 0 ]; then
    echo "ERROR: Build failed"
    exit 1
fi
echo "Done."
echo ""

# Step 4: Rename executable with version
echo "Step 4: Renaming executable with version..."
FINAL_NAME="Tyr-Desktop-${VERSION}-linux-amd64"
if [ -f "build/bin/Tyr-Desktop" ]; then
    mv build/bin/Tyr-Desktop "build/bin/${FINAL_NAME}"
    echo "Renamed to: ${FINAL_NAME}"
else
    echo "ERROR: Executable not found at build/bin/Tyr-Desktop"
    exit 1
fi
echo "Done."
echo ""

# Step 5: Updating system tray icon...
echo "Step 5: Updating system tray icon..."
if [ -f "build/linux/icon.png" ]; then
    cp build/linux/icon.png internal/resources/tyr.png
    echo "System tray icon updated"
fi
echo "Done."
echo ""

echo "========================================"
echo "Build completed successfully!"
echo "========================================"
echo "Executable: build/bin/${FINAL_NAME}"
echo "Version: ${VERSION}"
echo ""
