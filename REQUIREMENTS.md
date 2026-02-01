# Requirements - SpotiFLAC-CLI

This document outlines system requirements, runtime dependencies, and build requirements for SpotiFLAC-CLI.

---

## 🖥️ System Requirements

### End Users (Running Pre-Built Binaries)

#### Minimum Requirements

**All Platforms:**

- **Internet connection** - Required for downloading tracks
- **Storage**: 50 MB for binary + space for downloads (~30-120 MB per FLAC track)
- **Memory**: 100 MB RAM minimum

**Windows:**

- Windows 10 or later (AMD64 or ARM64)
- No additional dependencies required

**macOS:**

- macOS 10.13 High Sierra or later
- Intel (AMD64) or Apple Silicon (ARM64)
- No additional dependencies required

**Linux:**

- Any modern Linux distribution (Ubuntu 18.04+, Debian 10+, Fedora 30+, etc.)
- AMD64 or ARM64 architecture
- GLIBC 2.27+ (standard on modern Linux)
- No additional package dependencies

#### Recommended

- **CPU**: Any modern x86-64 or ARM64 processor
- **Memory**: 256 MB RAM for better performance
- **Storage**: 1 GB+ free space for multiple downloads
- **Network**: Broadband connection for faster downloads

### Special Platform Notes

#### Oracle Cloud Ampere A1

- Fully supported via **linux-arm64** build
- ARM64 architecture (AArch64)
- Tested on Oracle Linux 8/9 and Ubuntu 22.04
- No special configuration needed

#### Raspberry Pi

- **Raspberry Pi 4/5**: Use **linux-arm64** build
- **Raspberry Pi 3**: Use **linux-arm64** build (64-bit OS required)
- Requires 64-bit Raspberry Pi OS (Bullseye or later)
- Minimum 2 GB RAM recommended

#### Docker/Containers

- Any Linux-based container (Alpine, Ubuntu, Debian)
- Use **linux-amd64** or **linux-arm64** build
- No special runtime required

---

## 🔧 Runtime Dependencies

### FFmpeg (Bundled)

SpotiFLAC-CLI **includes FFmpeg functionality** internally via Go bindings. No external FFmpeg installation required.

**Internal handling:**

- Audio decoding and encoding
- Metadata embedding
- Cover art processing
- Format conversion

### Network Requirements

**Required endpoints:**

- `api.spotify.com` - Spotify metadata
- `song.link` - Cross-service URL resolution
- Tidal API endpoints
- Qobuz API endpoints
- Amazon Music API endpoints

**Firewall rules:**

- Allow outbound HTTPS (port 443)
- Allow outbound HTTP (port 80)

---

## 🏗️ Build Requirements

For developers building from source.

### Prerequisites

#### Required

1. **Go 1.21 or later**
   - Download: [golang.org/dl](https://golang.org/dl)
   - Verify: `go version`
   - Required for building Go code

2. **Git**
   - For cloning repository
   - Verify: `git --version`

3. **Build Tools**
   - **Windows**: No additional tools required
   - **Linux**: `build-essential` package (GCC, Make)
   - **macOS**: Xcode Command Line Tools

#### Recommended

- **Make** (optional) - For build automation
- **7-Zip** / `tar` - For creating release archives
- **UPX** (optional) - For additional binary compression

### Go Dependencies

All managed via `go.mod`:

```go
require (
    github.com/spf13/cobra v1.8.1           // CLI framework
    github.com/mattn/go-sqlite3 v1.14.24    // SQLite for history
    github.com/lib/pq v1.10.9               // PostgreSQL driver
)
```

**Installation:**

```bash
go mod download
# or
go mod tidy
```

### Platform-Specific Build Requirements

#### Windows

```powershell
# Install Go from golang.org
# No additional tools required

# Verify
go version
git --version
```

#### macOS

```bash
# Install Xcode Command Line Tools
xcode-select --install

# Install Go via Homebrew (optional)
brew install go

# Or download from golang.org

# Verify
go version
git --version
```

#### Linux (Ubuntu/Debian)

```bash
# Install Go and build tools
sudo apt update
sudo apt install -y golang-go git build-essential

# Verify
go version
git --version
gcc --version
```

#### Linux (Fedora/RHEL)

```bash
# Install Go and build tools
sudo dnf install -y golang git gcc

# Verify
go version
git --version
gcc --version
```

#### Oracle Linux (Ampere A1)

```bash
# Enable CodeReady Builder
sudo dnf install -y oracle-epel-release-el8
sudo dnf config-manager --enable ol8_codeready_builder

# Install Go and tools
sudo dnf install -y golang git gcc

# Verify
go version
```

---

## 🔨 Building

### Quick Build (Current Platform)

```bash
# Clone repository
git clone https://github.com/YOUR_USERNAME/SpotiFLAC-CLI.git
cd SpotiFLAC-CLI

# Build
go build -ldflags="-s -w" -o spotiflac main_cli.go

# Test
./spotiflac --version
```

### Cross-Compilation

Go supports cross-compilation out of the box:

```bash
# Windows AMD64 from Linux/Mac
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o spotiflac.exe main_cli.go

# Linux ARM64 from any platform
GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o spotiflac main_cli.go

# macOS Apple Silicon from any platform
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o spotiflac main_cli.go
```

### Multi-Platform Build

```bash
# Windows
.\scripts\build-all.bat

# Linux/macOS
chmod +x scripts/build-all.sh
./scripts/build-all.sh
```

This creates binaries for all 6 platforms:

- windows-amd64
- windows-arm64
- linux-amd64
- linux-arm64
- darwin-amd64
- darwin-arm64

---

## 📊 Build Specifications

### Build Flags

```bash
-ldflags="-s -w"
```

- `-s` - Strip symbol table (~15% size reduction)
- `-w` - Strip DWARF debug info (~15% size reduction)
- Combined: ~30% total size reduction

### Optional Optimization

```bash
# Additional size reduction with UPX
upx --best --lzma spotiflac

# Result: ~3-4 MB binary (vs ~9 MB uncompressed)
```

⚠️ **Note**: Some antivirus software flags UPX-compressed binaries as suspicious.

---

## 🧪 Testing Requirements

### Unit Testing

```bash
# Run all tests
go test ./...

# With coverage
go test -cover ./...

# Verbose output
go test -v ./...
```

### Integration Testing

```bash
# Test download functionality
./spotiflac https://open.spotify.com/track/... --auto

# Test JSON export
./spotiflac https://open.spotify.com/track/... -j

# Test batch download
./spotiflac https://open.spotify.com/album/... --auto -o ./test
```

---

## 🌐 Supported Architectures

| OS      | Architecture | GOOS      | GOARCH  | Notes                   |
| ------- | ------------ | --------- | ------- | ----------------------- |
| Windows | x86-64       | `windows` | `amd64` | 64-bit Windows          |
| Windows | ARM64        | `windows` | `arm64` | Surface Pro X, ARM PCs  |
| Linux   | x86-64       | `linux`   | `amd64` | Most Linux systems      |
| Linux   | ARM64        | `linux`   | `arm64` | Ampere A1, Raspberry Pi |
| macOS   | x86-64       | `darwin`  | `amd64` | Intel Macs              |
| macOS   | ARM64        | `darwin`  | `arm64` | M1/M2/M3/M4 Macs        |

---

## 🐛 Troubleshooting

### Build Issues

**"go: command not found"**

```bash
# Linux: Add to PATH
export PATH=$PATH:/usr/local/go/bin

# Windows: Add C:\Program Files\Go\bin to PATH
```

**"gcc: command not found" (Linux)**

```bash
# Ubuntu/Debian
sudo apt install build-essential

# Fedora/RHEL
sudo dnf install gcc
```

**CGO errors**

```bash
# Disable CGO if having issues
CGO_ENABLED=0 go build -ldflags="-s -w" -o spotiflac main_cli.go
```

### Runtime Issues

**"Permission denied" (Linux/Mac)**

```bash
chmod +x spotiflac
```

**"Cannot execute binary file"**

- Wrong architecture (e.g., ARM binary on x86 system)
- Download correct binary for your platform

**Network errors**

- Check internet connection
- Verify firewall allows HTTPS (port 443)
- Try with VPN if rate-limited

---

## 📦 Dependency Management

### Updating Dependencies

```bash
# Update all dependencies
go get -u ./...

# Update specific dependency
go get -u github.com/spf13/cobra@latest

# Clean up
go mod tidy
```

### Vendoring (Optional)

```bash
# Vendor all dependencies
go mod vendor

# Build with vendor
go build -mod=vendor -ldflags="-s -w" -o spotiflac main_cli.go
```

---

## 📝 Development Environment

### Recommended Tools

- **VS Code** with Go extension
- **GoLand** by JetBrains
- **Vim/Neovim** with vim-go
- **Terminal** with Go support

### Useful Go Tools

```bash
# Install useful tools
go install golang.org/x/tools/gopls@latest        # Language server
go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest  # Linter
go install golang.org/x/tools/cmd/goimports@latest  # Import formatter
```

---

## 🔐 Security Considerations

### Binary Integrity

- Verify checksums of downloaded binaries
- Build from source for maximum trust
- Inspect code before building

### Network Security

- Uses HTTPS for all API calls
- No telemetry or tracking
- No data sent to third parties except streaming services

---

## 📚 Additional Resources

- [Go Installation Guide](https://golang.org/doc/install)
- [Go Cross Compilation](https://go.dev/doc/install/source#environment)
- [Cobra CLI Framework](https://cobra.dev/)
- [SQLite Go Driver](https://github.com/mattn/go-sqlite3)

---

## ✅ Quick Setup Checklist

- [ ] Install Go 1.21+
- [ ] Install Git
- [ ] Clone repository
- [ ] Run `go mod download`
- [ ] Build with `go build -ldflags="-s -w" -o spotiflac main_cli.go`
- [ ] Test with `./spotiflac --version`
- [ ] Run sample download

**Need help?** Open an issue on GitHub.
