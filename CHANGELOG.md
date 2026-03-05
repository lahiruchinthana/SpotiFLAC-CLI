# Changelog

All notable changes to SpotiFLAC-CLI will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [1.1.4] - 2026-03-05

### 🐛 Fixed

#### Auto Fallback — Don't Abort When song.link Has No Tidal/Amazon Entry

- **`downloadWithAutoFallback`**: previously returned an error immediately if `GetAllURLsFromSpotify` failed (e.g. track not on Tidal or Amazon Music in song.link), causing Qobuz and Deezer to never be attempted. Now proceeds with empty `SongLinkURLs{}` and logs a debug warning, so all remaining services in the auto-order chain are still tried.

#### Qobuz — Always Set `lastErr` on ISRC Lookup Failure

- Qobuz block now properly sets `lastErr` at every failure stage:
  - Could not get Deezer URL for ISRC lookup
  - ISRC lookup itself failed
  - ISRC returned was empty
  - Previously, ISRC failures were silently swallowed, producing a misleading `"no services available"` error with no root cause.

#### Default `--auto-order` Updated

- Changed default from `tidal-amazon-qobuz` → `tidal-amazon-qobuz-deezer`
- Deezer is now included as the last-resort fallback for tracks not available on other services (Deezer lookups are Spotify-ID-based and don't depend on song.link).

#### Improved Terminal Error Message

- `"no services available"` → `"no configured services could find this track (auto-order: <value>)"` — now includes the actual service order that was tried.

---

## [1.1.3] - 2026-02-25

### ✨ Added - Deezer Service & Full v7.1.0 Backend Sync

Full alignment with upstream SpotiFLAC v7.1.0, adding Deezer as a fourth download
service and several backend improvements across the entire stack.

#### New Service: Deezer (via Yoinkify)

- **`backend/deezer.go`** — New Deezer downloader using the `yoinkify.lol` API
  - Downloads FLAC via POST to `yoinkify.lol/api/download` with `{url, format:"flac"}`
  - Full metadata pipeline: cover embedding, lyrics, genre, MusicBrainz, filename formatting
  - Integrated into both `--service deezer` and `--service auto` fallback chain

#### CLI Wiring for Deezer

- `--service` flag now accepts `deezer` as a valid value
- Long description updated: *"Get Spotify tracks in true FLAC from Tidal, Qobuz, Amazon Music & Deezer"*
- `downloadWithService()` — new `case "deezer"` block
- `downloadWithAutoFallback()` — Deezer added as fourth option in the auto chain

#### Backend Sync (v7.1.0)

- **`analysis.go`** — Added `GetMetadataWithFFprobe()` function; `AnalysisResult` now includes `Bitrate` field
- **`cover.go`** — Fixed `{date}` template variable support in cover file naming
- **`filename.go`** — Fixed `{date}` template variable support in filename formatting
- **`lyrics.go`** — Spotify Lyrics API added as first-priority source before Musixmatch/NetEase
- **`uploader.go`** — Dynamic upload URL discovery replaces hardcoded endpoint
- **`musicbrainz.go`** — New MusicBrainz integration for genre tagging via `--embed-genre`
- **`metadata.go`** — `Genre` field added to `TrackMetadata` struct
- **`tidal.go`** — API endpoint cleanup; genre field populated from metadata
- **`qobuz.go`** — API endpoint cleanup; genre field populated from metadata
- **`amazon.go`** — Genre field populated from metadata

#### Preserved CLI-Specific Behaviour

- `spotfetch.go` — Custom TOTP obfuscation (byte-array XOR, hex+base32 encoding) retained; upstream plain-string secret not adopted
- All new downloader `Description` fields set to `"https://downbot.app"` instead of upstream GitHub URL

---

## [1.1.2] - 2026-02-18

### 🐛 Fixed

- **Critical:** Fixed playlist downloads failing with "Unknown content type" error
  - Resolved type mismatch between `PlaylistResponsePayload` (value) and `*PlaylistResponsePayload` (pointer)
  - Playlists now download correctly with all metadata and tracks

---

## [1.1.1] - 2026-02-14

### ✨ Added - Enhanced Progress Display & Command-Line Options

Enhanced command-line interface with additional options and improved progress display!

#### New Features

- **Enhanced progress display** - Download progress now displays with percentage, speed, and ETA
  - Format: `[download] 8.8% of 26.19MiB at 2.98MiB/s ETA 00:08`
  - Real-time speed calculation and accurate time estimates
- **New command-line flags:**
  - `--print-json` - Print metadata as JSON (alias for `--dump-json`)
  - `--no-warnings` - Suppress warning messages
  - `--newline` - Output progress on new lines instead of overwriting
  - `--progress` - Show/hide download progress (default: enabled)
  - `--print <field>` - Print specific metadata field (title, artist, album, url, etc.)

#### New Command-Line Options

```bash
--print-json              Print metadata as JSON
--no-warnings            Suppress warning messages
--newline                Output progress on new lines
--progress               Show download progress (default: true)
--print <field>          Print specific field from metadata
```

#### Quick Examples

```bash
# Download with progress on new lines
spotiflac <URL> --auto --newline

# Print only the track title
spotiflac <URL> --print title

# Print metadata as JSON without downloading
spotiflac <URL> --print-json

# Download silently (no warnings or progress)
spotiflac <URL> --auto --no-warnings --progress=false --quiet
```

### 🐛 Bug Fixes

- **Fixed MP3 subdirectory creation** - MP3 files are now saved in the current directory instead of creating `MP3/` subdirectory
- **Fixed MP3 conversion for existing files** - Existing FLAC files are now properly converted when using `--output-format mp3`
- **Improved quiet mode** - Better handling of quiet/verbose flags with new `--no-warnings` option

### 🔧 Technical Details

- Progress display respects `--newline` flag for continuous output
- `--print` supports common fields: title, name, artist, artists, album, album_name, album_artist, release_date, date, spotify_id, id, duration, duration_ms, track_number, url, spotify_url
- Progress calculation includes total file size, download speed, and ETA
- Warning messages can be suppressed independently from info messages

---

## [1.1.0] - 2026-02-13

### ✨ Added - MP3 Format Support

SpotiFLAC-CLI now supports automatic conversion to MP3 format after downloading lossless FLAC files!

#### New Features

- **MP3 conversion** - Convert downloaded FLAC to MP3 automatically with `--output-format mp3`
- **Configurable bitrate** - Set MP3 bitrate with `--mp3-bitrate` (128k, 192k, 256k, 320k)
- **M4A support** - Also supports M4A format conversion with `--output-format m4a`
- **Automatic metadata** - All metadata (tags, cover art, lyrics) preserved during conversion
- **Clean workflow** - Original FLAC file automatically removed after successful conversion

#### New Command-Line Options

```bash
--output-format <format>   Output format: flac (default), mp3, m4a
--mp3-bitrate <bitrate>    MP3 bitrate (default: 320k)
                          Options: 128k, 192k, 256k, 320k
```

#### Quick Examples

```bash
# Download as MP3 320kbps (highest quality)
spotiflac <URL> --auto --output-format mp3

# Download as MP3 256kbps
spotiflac <URL> --auto --output-format mp3 --mp3-bitrate 256k

# Download album as MP3
spotiflac <ALBUM_URL> --auto --output-format mp3
```

#### Technical Details

- Uses FFmpeg with libmp3lame encoder
- ID3v2.3 tags for maximum compatibility
- Preserves all metadata: title, artist, album, year, track numbers, ISRC, lyrics, cover art
- Organized output: MP3 files saved in `MP3/` subdirectory

### 🔄 Updated to align with SpotiFLAC v7.0.9

This release brings SpotiFLAC-CLI up to date with the main SpotiFLAC repository (v7.0.7 → v7.0.9), incorporating critical bug fixes and new features.

### 🐛 Fixed

- **Fixed missing artist and publisher metadata** - Albums now properly include publisher/label information (#473)
- **Fixed lyric embedding failures** - Lyrics now embed correctly into FLAC files
- **Fixed Qobuz API endpoint** - Updated to new Jumo API endpoint with proper headers
- **Fixed Amazon Music API** - Updated to new API structure with decryption support
- **Fixed incorrect disc numbers** - Multi-disc albums now show correct disc numbers in artist discography
- **Fixed ISRC metadata** - ISRC codes are now properly embedded (was incorrectly using Spotify track IDs)

### ✨ Added

- **ISRC metadata embedding** - Proper ISRC codes now embedded in downloaded tracks for better identification (#472)
- **Amazon Music decryption support** - Automatic FFmpeg decryption for DRM-protected Amazon streams
- **Codec detection** - Automatic detection and handling of FLAC/M4A codecs from Amazon Music
- **User-Agent headers** - Added proper headers for improved API compatibility
- **Enhanced error handling** - Better error messages for API failures

### 🔧 Changed

- Updated Qobuz API from `/file` to `/get` endpoint
- Improved Amazon Music download flow with decryption pipeline
- Enhanced metadata embedding with ISRC field
- Better disc information handling for multi-disc albums

### 📦 Backend Changes (from SpotiFLAC v7.0.8 & v7.0.9)

- Qobuz and Amazon Music API fixes
- Metadata tagging failure fixes
- Year issues in lyrics and cover metadata resolved
- Automatic ffmpeg detection improvements
- SpotFetch API fallback suggestions

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
