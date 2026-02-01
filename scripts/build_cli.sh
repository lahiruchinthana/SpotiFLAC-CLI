#!/bin/bash

# Change to project root
cd "$(dirname "$0")/.."

echo "========================================="
echo "   SpotiFLAC CLI Build Script"
echo "========================================="

# Exit on error
set -e

# Clean old binary
rm -f spotiflac

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

# Build for current platform
echo ""
echo "Building for current platform..."
go build -ldflags="-s -w" -o spotiflac main_cli.go

# Make executable
chmod +x spotiflac

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
echo "Binary: ./spotiflac"
echo "Run: ./spotiflac --help"
echo ""
