#!/bin/bash

# Change to project root
cd "$(dirname "$0")/.."

echo "========================================="
echo "   SpotiFLAC CLI - Multi-Platform Build"
echo "========================================="
echo ""

# Set version
VERSION="7.0.7"

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Clean build directory
echo "Cleaning build directory..."
rm -rf build
mkdir -p build/{windows-amd64,windows-arm64,linux-amd64,linux-arm64,darwin-amd64,darwin-arm64}

# Backup original go.mod
if [ -f "go.mod" ] && [ ! -f "go.mod.backup" ]; then
    echo "Backing up go.mod..."
    cp go.mod go.mod.backup
fi

# Use CLI go.mod
echo "Switching to CLI dependencies..."
cp go_cli.mod go.mod

# Install dependencies
echo "Installing dependencies..."
go mod tidy

echo ""
echo "Building for multiple platforms..."
echo ""

# Build function
build() {
    local os=$1
    local arch=$2
    local num=$3
    local total=$4
    local ext=$5
    
    echo -e "${BLUE}[$num/$total]${NC} Building $os $arch..."
    
    GOOS=$os GOARCH=$arch go build \
        -ldflags="-s -w -X main.version=$VERSION" \
        -o "build/$os-$arch/spotiflac$ext" \
        main_cli.go
    
    if [ $? -eq 0 ]; then
        echo -e "  ${GREEN}✓${NC} build/$os-$arch/spotiflac$ext"
        return 0
    else
        echo -e "  ${RED}✗${NC} ERROR: $os $arch build failed"
        return 1
    fi
}

# Build for all platforms
build windows amd64 1 6 .exe || exit 1
build windows arm64 2 6 .exe || exit 1
build linux amd64 3 6 "" || exit 1
build linux arm64 4 6 "" || exit 1
build darwin amd64 5 6 "" || exit 1
build darwin arm64 6 6 "" || exit 1

echo ""
echo "Copying documentation to build directories..."
for dir in build/*/; do
    cp docs/README_CLI.md "$dir/README.md" 2>/dev/null || true
    cp LICENSE "$dir/LICENSE" 2>/dev/null || true
done

# Make Linux and macOS binaries executable
chmod +x build/linux-*/* 2>/dev/null || true
chmod +x build/darwin-*/* 2>/dev/null || true

echo ""
echo "========================================="
echo "   Build Summary"
echo "========================================="
echo ""
echo "Platforms built:"
echo "  • Windows AMD64 (x64)"
echo "  • Windows ARM64 (ARM)"
echo "  • Linux AMD64 (x64)"
echo "  • Linux ARM64 (ARM)"
echo "  • macOS AMD64 (Intel)"
echo "  • macOS ARM64 (Apple Silicon)"
echo ""
echo "Output directory: build/"
echo ""
ls -1 build/
echo ""

# Get file sizes
echo "Build sizes:"
for dir in build/*/; do
    platform=$(basename "$dir")
    file=$(find "$dir" -name "spotiflac*" -type f)
    if [ -n "$file" ]; then
        size=$(du -h "$file" | cut -f1)
        printf "  %-20s %s\n" "$platform:" "$size"
    fi
done

# Restore original go.mod
echo ""
echo "Restoring original go.mod..."
if [ -f "go.mod.backup" ]; then
    mv go.mod.backup go.mod
fi

echo ""
echo "========================================="
echo "   Build Complete!"
echo "========================================="
echo ""
echo "To test: ./build/$(uname -s | tr '[:upper:]' '[:lower:]')-amd64/spotiflac --help"
echo ""
