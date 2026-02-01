# SpotiFLAC-CLI

[![GitHub Release](https://img.shields.io/github/v/release/YOUR_USERNAME/SpotiFLAC-CLI?style=for-the-badge)](https://github.com/YOUR_USERNAME/SpotiFLAC-CLI/releases)
[![Downloads](https://img.shields.io/github/downloads/YOUR_USERNAME/SpotiFLAC-CLI/total?style=for-the-badge)](https://github.com/YOUR_USERNAME/SpotiFLAC-CLI/releases)
[![License](https://img.shields.io/github/license/YOUR_USERNAME/SpotiFLAC-CLI?style=for-the-badge)](LICENSE)

**Command-line tool to download Spotify tracks in lossless FLAC format from Tidal, Qobuz & Amazon Music — no account required.**

```bash
spotiflac https://open.spotify.com/track/... --auto
```

---

## ✨ Features

- 🎵 **True Lossless Audio** - Download 16-bit/24-bit FLAC from official sources
- 🔄 **Multi-Service Support** - Tidal, Qobuz, and Amazon Music
- 🎯 **Smart Auto Mode** - Automatically tries multiple services with fallback
- 📝 **Rich Metadata** - Album art, lyrics, copyright, and ISRC codes
- 📦 **Batch Downloads** - Albums, playlists, and multiple tracks
- 🌍 **Cross-Platform** - Windows, macOS, Linux (AMD64 & ARM64)
- 🚀 **Single Binary** - No dependencies, just download and run
- 📋 **JSON Export** - Get metadata and download links without downloading
- 🔍 **ISRC Matching** - Precise track matching across services
- ⚡ **Fast & Lightweight** - Built with Go, ~9-10 MB binary

---

## 📥 Installation

### Download Pre-Built Binaries

Download the latest release for your platform:

#### Windows

```powershell
# Download windows-amd64.zip or windows-arm64.zip
# Extract and run spotiflac.exe
```

#### macOS

```bash
# Intel Macs
curl -LO https://github.com/YOUR_USERNAME/SpotiFLAC-CLI/releases/latest/download/spotiflac-darwin-amd64.tar.gz
tar -xzf spotiflac-darwin-amd64.tar.gz
chmod +x spotiflac
sudo mv spotiflac /usr/local/bin/

# Apple Silicon (M1/M2/M3/M4)
curl -LO https://github.com/YOUR_USERNAME/SpotiFLAC-CLI/releases/latest/download/spotiflac-darwin-arm64.tar.gz
tar -xzf spotiflac-darwin-arm64.tar.gz
chmod +x spotiflac
sudo mv spotiflac /usr/local/bin/
```

#### Linux

```bash
# AMD64 (x86-64)
wget https://github.com/YOUR_USERNAME/SpotiFLAC-CLI/releases/latest/download/spotiflac-linux-amd64.tar.gz
tar -xzf spotiflac-linux-amd64.tar.gz
chmod +x spotiflac
sudo mv spotiflac /usr/local/bin/

# ARM64 (Oracle Ampere A1, Raspberry Pi 4/5)
wget https://github.com/YOUR_USERNAME/SpotiFLAC-CLI/releases/latest/download/spotiflac-linux-arm64.tar.gz
tar -xzf spotiflac-linux-arm64.tar.gz
chmod +x spotiflac
sudo mv spotiflac /usr/local/bin/
```

### Build from Source

```bash
git clone https://github.com/YOUR_USERNAME/SpotiFLAC-CLI.git
cd SpotiFLAC-CLI
go build -ldflags="-s -w" -o spotiflac main_cli.go
```

**[→ See REQUIREMENTS.md for build dependencies](REQUIREMENTS.md)**

---

## 🚀 Quick Start

```bash
# Simple download (tries Tidal first)
spotiflac https://open.spotify.com/track/3n3Ppam7vgaVa1iaRUc9Lp

# Auto mode - tries all services (recommended)
spotiflac https://open.spotify.com/track/3n3Ppam7vgaVa1iaRUc9Lp --auto

# High quality with lyrics and cover art
spotiflac https://open.spotify.com/track/3n3Ppam7vgaVa1iaRUc9Lp --auto --auto-quality 24 --embed-lyrics --max-quality-cover

# Download entire album
spotiflac https://open.spotify.com/album/4a6NzYL1YHRUgx9e3YZI6I --auto -o ./Music

# Download playlist with delay between tracks
spotiflac https://open.spotify.com/playlist/37i9dQZF1DXcBWIGoYBM5M --auto --delay 2.0

# Get metadata only (no download)
spotiflac https://open.spotify.com/track/3n3Ppam7vgaVa1iaRUc9Lp -j
```

---

## 📖 Usage

### Basic Commands

```bash
spotiflac [URL] [OPTIONS]
```

### Common Options

| Option                | Short | Description                           | Example      |
| --------------------- | ----- | ------------------------------------- | ------------ |
| `--auto`              |       | Auto mode - tries multiple services   | `--auto`     |
| `--service`           | `-s`  | Specific service (tidal/qobuz/amazon) | `-s tidal`   |
| `--quality`           | `-q`  | Quality setting                       | `-q HI_RES`  |
| `--output`            | `-o`  | Output directory                      | `-o ./Music` |
| `--dump-json`         | `-j`  | Print metadata as JSON and exit       | `-j`         |
| `--embed-lyrics`      | `-l`  | Embed lyrics                          | `-l`         |
| `--max-quality-cover` | `-c`  | High quality cover art                | `-c`         |
| `--delay`             | `-d`  | Delay between downloads (seconds)     | `-d 2.0`     |
| `--verbose`           | `-v`  | Verbose output                        | `-v`         |
| `--quiet`             |       | Minimal output                        | `--quiet`    |

### Quality Settings

**Tidal:**

- `LOSSLESS` - 16-bit FLAC (CD quality, ~30-40 MB/track)
- `HI_RES` - 24-bit FLAC/MQA (~80-120 MB/track)

**Qobuz:**

- `6` - 16-bit FLAC (CD quality)
- `7` - 24-bit FLAC (Hi-Res)
- `27` - Hi-Res+ (highest available)

**Amazon:**

- Auto-selects best available quality

### Auto Mode Examples

```bash
# Auto mode with 24-bit quality preference
spotiflac URL --auto --auto-quality 24

# Custom service order
spotiflac URL --auto --auto-order qobuz-tidal-amazon

# Auto mode with all enhancements
spotiflac URL --auto --auto-quality 24 -l -c
```

### Batch Downloads

```bash
# Album
spotiflac https://open.spotify.com/album/... --auto

# Playlist
spotiflac https://open.spotify.com/playlist/... --auto --delay 2.0

# Multiple tracks from file
cat urls.txt | while read url; do
    spotiflac "$url" --auto --delay 1.0
done
```

### Metadata Export

```bash
# Get metadata as JSON without downloading
spotiflac URL --dump-json

# Get metadata as JSON with pretty formatting
spotiflac URL -j | jq '.'

# Save metadata and download
spotiflac URL --write-info-json --auto

# Extract only download links
spotiflac URL -j | jq '.download_links'
```

---

## 🎯 Examples

### Single Track (Best Quality)

```bash
spotiflac https://open.spotify.com/track/1k0JAiH11gHL9dc5dfQjQr \
  --auto \
  --auto-quality 24 \
  --embed-lyrics \
  --max-quality-cover \
  -o ~/Music
```

### Album with Lyrics

```bash
spotiflac https://open.spotify.com/album/4a6NzYL1YHRUgx9e3YZI6I \
  --auto \
  --embed-lyrics \
  --track-number \
  -o ~/Music
```

### Playlist with Rate Limiting

```bash
spotiflac https://open.spotify.com/playlist/37i9dQZF1DXcBWIGoYBM5M \
  --auto \
  --delay 2.0 \
  -o ~/Playlists
```

### Check Availability

```bash
# Get metadata to see which services have the track
spotiflac https://open.spotify.com/track/... -j | jq '.available_services'
```

---

## 🏗️ Supported Platforms

| Platform    | Architecture     | Binary Size | Notes                                  |
| ----------- | ---------------- | ----------- | -------------------------------------- |
| **Windows** | AMD64 (x64)      | ~9.7 MB     | Windows 10/11                          |
| **Windows** | ARM64            | ~8.9 MB     | Surface Pro X, ARM PCs                 |
| **Linux**   | AMD64 (x64)      | ~9.3 MB     | Ubuntu, Debian, Fedora                 |
| **Linux**   | ARM64            | ~8.8 MB     | **Oracle Ampere A1**, Raspberry Pi 4/5 |
| **macOS**   | AMD64 (Intel)    | ~9.5 MB     | Intel Macs                             |
| **macOS**   | ARM64 (M1/M2/M3) | ~8.9 MB     | Apple Silicon                          |

---

## 🔧 Configuration

SpotiFLAC-CLI uses command-line flags for all configuration. For convenience, you can create shell aliases:

### Bash/Zsh (~/.bashrc or ~/.zshrc)

```bash
# Quick download with auto mode
alias sfdl='spotiflac --auto --auto-quality 24'

# Download with all metadata
alias sfdl-full='spotiflac --auto --auto-quality 24 --embed-lyrics --max-quality-cover'

# Metadata only
alias sfmeta='spotiflac --dump-json'
```

### PowerShell (Profile)

```powershell
function sfdl { spotiflac $args --auto --auto-quality 24 }
function sfdl-full { spotiflac $args --auto --auto-quality 24 --embed-lyrics --max-quality-cover }
function sfmeta { spotiflac $args --dump-json }
```

---

## 📋 JSON Metadata Format

```json
{
  "available_services": ["tidal", "amazon", "qobuz"],
  "download_links": {
    "tidal": {
      "url": "https://listen.tidal.com/track/...",
      "quality": "LOSSLESS (16-bit) / HI_RES (24-bit)"
    },
    "amazon": {
      "url": "https://music.amazon.com/...",
      "quality": "Variable (auto-select)"
    },
    "qobuz": {
      "quality": "6 (16-bit) / 7 (24-bit) / 27 (Hi-Res)",
      "note": "Available via ISRC lookup"
    }
  },
  "metadata": {
    "track": {
      "name": "Track Name",
      "artists": "Artist Name",
      "album_name": "Album Name",
      "duration_ms": 180000,
      "release_date": "2025-10-03",
      "isrc": "USRC12345678",
      "plays": "123456789"
    }
  }
}
```

---

## ❓ FAQ

### Is this free?

Yes, completely free and open source. No account or subscription required.

### Can this get my Spotify account banned?

No. This tool has no connection to your Spotify account. It uses Spotify's public metadata through web player reverse engineering, not user authentication.

### Where does the audio come from?

Audio is downloaded from Tidal, Qobuz, and Amazon Music via third-party APIs. The tool matches Spotify tracks using ISRC codes.

### Why does metadata fetching sometimes fail?

Usually due to IP rate limiting. Wait and try again, or use a VPN. The `--auto` mode helps by trying multiple services.

### Does this work on servers/headless systems?

Yes! SpotiFLAC-CLI is perfect for servers, Docker containers, and headless systems. No GUI required.

### What about the original SpotiFLAC?

This is a CLI port of [afkarxyz/SpotiFLAC](https://github.com/afkarxyz/SpotiFLLAC). The original has a GUI (Wails/React). This version focuses on command-line usage for automation, servers, and power users.

---

## 🤝 Credits & Attribution

**This project is a command-line port of the original [SpotiFLAC](https://github.com/afkarxyz/SpotiFLAC) by [@afkarxyz](https://github.com/afkarxyz).**

All core download functionality, service integrations, and metadata handling come from the original project. This CLI version:

- Removes the Wails GUI and React frontend
- Adds Cobra CLI framework for command-line interface
- Implements JSON metadata export
- Optimizes for automation and scripting
- Maintains 100% feature parity with backend logic

**Original Project:** https://github.com/afkarxyz/SpotiFLAC

### API Credits

- **Tidal**: [hifi-api](https://github.com/binimum/hifi-api)
- **Qobuz**: [dabmusic.xyz](https://dabmusic.xyz), [squid.wtf](https://squid.wtf), [jumo-dl](https://jumo-dl.pages.dev/)
- **Song Matching**: [song.link](https://song.link) API

---

## 📜 Disclaimer

This project is for **educational and private use only**. The developer does not condone or encourage copyright infringement.

**SpotiFLAC-CLI** is a third-party tool and is not affiliated with, endorsed by, or connected to Spotify, Tidal, Qobuz, Amazon Music, or any other streaming service.

You are solely responsible for:

1. Ensuring your use complies with local laws
2. Reading and adhering to Terms of Service of respective platforms
3. Any legal consequences from misuse of this tool

The software is provided "as is", without warranty of any kind. The author assumes no liability for any bans, damages, or legal issues arising from its use.

---

## 📄 License

This project inherits its license from the original [SpotiFLAC](https://github.com/afkarxyz/SpotiFLAC) project. See [LICENSE](LICENSE) for details.

---

## 🛠️ Development

See [RELEASE_GUIDE.md](RELEASE_GUIDE.md) for build and release instructions.

**[→ Build Requirements](REQUIREMENTS.md)** | **[→ Release Guide](RELEASE_GUIDE.md)** | **[→ Changelog](CHANGELOG.md)**

---

<div align="center">

**Made with ❤️ for the music community**

⭐ Star this repo if you find it useful!

</div>
