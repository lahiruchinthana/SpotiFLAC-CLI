# GitHub Release Creation Guide - SpotiFLAC-CLI v1.0.0

This guide walks you through creating your first GitHub release.

---

## 📦 Release Files Ready

All 6 binary archives have been created in the `build/` directory:
- `spotiflac-windows-amd64.zip` (9.66 MB)
- `spotiflac-windows-arm64.zip` (8.89 MB)
- `spotiflac-linux-amd64.tar.gz` (9.33 MB)
- `spotiflac-linux-arm64.tar.gz` (8.75 MB)
- `spotiflac-darwin-amd64.tar.gz` (9.52 MB)
- `spotiflac-darwin-arm64.tar.gz` (8.94 MB)

---

## 🚀 Step-by-Step Release Process

### Step 1: Navigate to Releases
1. Go to https://github.com/lahiruchinthana/SpotiFLAC-CLI
2. Click **"Releases"** in the right sidebar
3. Click **"Draft a new release"** button

### Step 2: Choose a Tag
- Click the **"Choose a tag"** dropdown
- Select **`v1.0.0`** (existing tag)
- It should show as "Existing tag"

### Step 3: Release Title
Enter this exactly:
```
v1.0.0 - Initial CLI Release
```

### Step 4: Release Description
Copy and paste this entire markdown:

```markdown
# SpotiFLAC-CLI v1.0.0 🎉

**Command-line tool to download Spotify tracks in lossless FLAC format from Tidal, Qobuz & Amazon Music.**

## ✨ What's New

This is the initial release of SpotiFLAC-CLI - a command-line port of the original [SpotiFLAC](https://github.com/afkarxyz/SpotiFLAC) GUI application.

### Core Features
- 🎵 **Multi-service downloads** - Tidal, Qobuz, Amazon Music
- 🎯 **Smart auto mode** - Automatic service fallback
- 📋 **JSON metadata export** - Get track info without downloading (`--dump-json`, `-j`)
- 📦 **Batch downloads** - Albums and playlists with rate limiting
- 🌍 **Cross-platform** - Windows, macOS, Linux (AMD64 & ARM64)
- 🚀 **Single binary** - No dependencies, ~9-10 MB per platform

### Quality Options
- **Tidal**: LOSSLESS (16-bit) / HI_RES (24-bit)
- **Qobuz**: Quality 6 (16-bit) / 7 (24-bit) / 27 (Hi-Res)
- **Amazon Music**: Auto-select best available

### Platform Support
- ✅ Windows 10/11 (AMD64 & ARM64)
- ✅ macOS Intel & Apple Silicon
- ✅ Linux AMD64 & ARM64 (Oracle Ampere A1, Raspberry Pi 4/5)

## 📥 Installation

Download the appropriate binary for your platform below, extract, and run!

### Quick Start
```bash
# Simple download
spotiflac https://open.spotify.com/track/... --auto

# High quality with metadata
spotiflac URL --auto --auto-quality 24 --embed-lyrics --max-quality-cover

# Get metadata only (no download)
spotiflac URL -j
```

## 📊 Binary Sizes

| Platform | Architecture | Size |
|----------|-------------|------|
| Windows | AMD64 | 9.66 MB |
| Windows | ARM64 | 8.89 MB |
| Linux | AMD64 | 9.33 MB |
| Linux | ARM64 | 8.75 MB |
| macOS | Intel | 9.52 MB |
| macOS | Apple Silicon | 8.94 MB |

## 🤝 Credits

Based on [SpotiFLAC v7.0.7](https://github.com/afkarxyz/SpotiFLAC) by [@afkarxyz](https://github.com/afkarxyz).

All download functionality preserved from original project. CLI additions: Cobra framework, JSON export, batch processing improvements, and platform optimization.

## 📖 Documentation

- [README](https://github.com/lahiruchinthana/SpotiFLAC-CLI#readme) - Full usage guide
- [CHANGELOG](https://github.com/lahiruchinthana/SpotiFLAC-CLI/blob/main/CHANGELOG.md) - Detailed version history
- [Build from Source](https://github.com/lahiruchinthana/SpotiFLAC-CLI/blob/main/REQUIREMENTS.md) - Requirements for building

---

**Full Changelog**: https://github.com/lahiruchinthana/SpotiFLAC-CLI/blob/main/CHANGELOG.md
```

### Step 5: Upload Binaries
1. Scroll down to **"Attach binaries by dropping them here or selecting them"**
2. Click to browse or drag these 6 files from `build/` directory:
   - ✅ `spotiflac-windows-amd64.zip`
   - ✅ `spotiflac-windows-arm64.zip`
   - ✅ `spotiflac-linux-amd64.tar.gz`
   - ✅ `spotiflac-linux-arm64.tar.gz`
   - ✅ `spotiflac-darwin-amd64.tar.gz`
   - ✅ `spotiflac-darwin-arm64.tar.gz`
3. Wait for all uploads to complete (you'll see progress bars)

### Step 6: Release Settings
- ✅ **Set as the latest release** - Leave CHECKED
- ❌ **Set as a pre-release** - Leave UNCHECKED
- ❌ **Create a discussion** - Optional, can skip

### Step 7: Publish
Click **"Publish release"** button!

---

## ✅ Pre-Publish Checklist

Before clicking "Publish release", verify:

- [ ] Tag: `v1.0.0` selected
- [ ] Title: `v1.0.0 - Initial CLI Release`
- [ ] Description: Markdown copied completely
- [ ] 6 binary files uploaded (2 ZIP + 4 TAR.GZ)
- [ ] "Set as latest release" is checked
- [ ] "Set as pre-release" is unchecked

---

## 🎉 After Publishing

Your release will be live at:
https://github.com/lahiruchinthana/SpotiFLAC-CLI/releases/tag/v1.0.0

Users can download binaries using:
```
https://github.com/lahiruchinthana/SpotiFLAC-CLI/releases/latest/download/spotiflac-[platform]-[arch].[zip|tar.gz]
```

The README badges will automatically update to show v1.0.0!

---

## 📝 Notes

- GitHub automatically creates source code archives (`.zip` and `.tar.gz`)
- The release appears on your repo's main page
- Download links in README will work immediately
- You can edit the release later if needed

---

**Ready to go! Just follow the steps above and upload the 6 files from `build/` directory.**
