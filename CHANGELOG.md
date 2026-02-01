# Changelog

All notable changes to SpotiFLAC-CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.0.0] - 2026-02-01

### 🎉 Initial Release

First public release of SpotiFLAC-CLI - a command-line port of the original [SpotiFLAC](https://github.com/afkarxyz/SpotiFLAC) GUI application.

### ✨ Added

#### Core Features

- **Command-line interface** powered by Cobra framework
- **Multi-service downloads** from Tidal, Qobuz, and Amazon Music
- **Smart auto mode** with automatic service fallback
- **Batch downloads** for albums and playlists
- **ISRC-based track matching** for precise identification
- **FLAC format** with 16-bit/24-bit quality options
- **Rich metadata embedding** including artist, album, track info, ISRC, and copyright
- **Cover art embedding** with optional high-quality mode
- **Lyrics embedding** from multiple sources

#### JSON Metadata Export

- `--dump-json` / `-j` flag to output metadata as JSON without downloading
- `--write-info-json` flag to save metadata alongside downloads
- JSON output includes:
  - Complete track/album/playlist metadata
  - Download links for Tidal, Amazon Music, and Qobuz
  - Available services list
  - Quality information per service
  - Usage examples for CLI commands

#### Quality Control

- Service-specific quality settings (LOSSLESS/HI_RES for Tidal, 6/7/27 for Qobuz)
- `--auto-quality` flag for intelligent quality selection (16/24-bit)
- Automatic quality fallback when preferred quality unavailable

#### Batch Processing

- Album downloads with full track metadata
- Playlist downloads with rate limiting support
- `--delay` flag for controlling download intervals
- Progress tracking for multi-track downloads

#### Output Customization

- `--output` / `-o` flag for custom output directory
- `--filename-format` for template-based naming
- `--track-number` flag to include track numbers in filenames
- Automatic file organization

#### Platform Support

- **Windows**: AMD64 (x64) and ARM64
- **Linux**: AMD64 (x64) and ARM64 (Oracle Ampere A1, Raspberry Pi compatible)
- **macOS**: Intel (AMD64) and Apple Silicon (ARM64)
- Single-binary distribution (~8-10 MB per platform)
- No runtime dependencies except FFmpeg (handled internally)

#### CLI Features

- `--verbose` / `-v` flag for detailed logging
- `--quiet` flag for minimal output
- `--version` to display version information
- `--help` for comprehensive usage instructions
- Exit codes for scripting integration
- Shell-friendly output formatting

### 🏗️ Infrastructure

#### Build System

- Multi-platform build scripts (`build-all.sh`, `build-all.bat`)
- Single-platform build scripts (`build_cli.sh`, `build_cli.bat`)
- Go module management with minimal dependencies
- Cross-compilation support for all 6 platforms
- Automated binary stripping and optimization

#### Documentation

- Comprehensive README.md with examples
- REQUIREMENTS.md for system and build requirements
- RELEASE_GUIDE.md for maintainers
- LICENSE file (inherited from original project)

### 📦 Technical Details

#### Dependencies

- **Go 1.21+** - Primary language
- **Cobra v1.8.1** - CLI framework
- **SQLite v1.14.24** - Download history (mattn/go-sqlite3)
- **lib/pq v1.10.9** - Database driver
- All backend dependencies from original SpotiFLAC

#### Architecture

- **Backend**: 100% preserved from original SpotiFLAC v7.0.7
  - `amazon.go`, `tidal.go`, `qobuz.go` - Service downloaders
  - `spotify_metadata.go` - Spotify API integration
  - `songlink.go` - Cross-service URL resolution
  - `metadata.go`, `cover.go`, `lyrics.go` - Metadata handling
  - `ffmpeg.go` - Audio processing
  - `history.go` - Download tracking
- **Frontend**: Removed Wails GUI and React components
- **CLI**: New `main_cli.go` with Cobra framework

### 🎯 Use Cases

Perfect for:

- **Automation** - Scripting and batch operations
- **Servers** - Headless/remote systems without GUI
- **CI/CD** - Integration into build pipelines
- **Power Users** - Fast command-line workflows
- **Docker** - Containerized music downloading
- **Cloud** - Oracle Ampere A1, AWS Graviton instances

### 🤝 Attribution

Based on [SpotiFLAC v7.0.7](https://github.com/afkarxyz/SpotiFLAC) by [@afkarxyz](https://github.com/afkarxyz).

All download functionality, service integrations, and metadata handling preserved from original project. CLI additions:

- Cobra command framework
- JSON export functionality
- Batch processing improvements
- Platform optimization

### 📊 Binary Sizes

| Platform | Architecture | Size    |
| -------- | ------------ | ------- |
| Windows  | AMD64        | 9.66 MB |
| Windows  | ARM64        | 8.89 MB |
| Linux    | AMD64        | 9.33 MB |
| Linux    | ARM64        | 8.75 MB |
| macOS    | AMD64        | 9.52 MB |
| macOS    | ARM64        | 8.94 MB |

### 🔗 Links

- **Original Project**: https://github.com/afkarxyz/SpotiFLAC
- **Tidal API**: [hifi-api](https://github.com/binimum/hifi-api)
- **Qobuz Sources**: [dabmusic.xyz](https://dabmusic.xyz), [squid.wtf](https://squid.wtf), [jumo-dl](https://jumo-dl.pages.dev/)
- **Song Matching**: [song.link](https://song.link)

---

## [Unreleased]

### Planned for v1.1.0

- Shell completion scripts (bash, zsh, fish)
- Docker image for easy deployment
- Configuration file support (~/.spotiflac.yaml)
- Parallel downloads for playlists
- Resume incomplete downloads
- Progress bars using terminal UI library

---

## Version History

- **v1.0.0** (2026-02-01) - Initial CLI release
- Based on SpotiFLAC v7.0.7 backend

---

**Format**: [Keep a Changelog](https://keepachangelog.com/)  
**Versioning**: [Semantic Versioning](https://semver.org/)
