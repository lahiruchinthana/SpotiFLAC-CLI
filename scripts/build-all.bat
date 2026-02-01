@echo off
setlocal enabledelayedexpansion

REM Change to project root
cd /d "%~dp0\.."

echo =========================================
echo    SpotiFLAC CLI - Multi-Platform Build
echo =========================================
echo.

REM Set version
set VERSION=7.0.7

REM Clean build directory
echo Cleaning build directory...
if exist build (
    rd /s /q build
)
mkdir build\windows-amd64
mkdir build\windows-arm64
mkdir build\linux-amd64
mkdir build\linux-arm64
mkdir build\darwin-amd64
mkdir build\darwin-arm64

REM Backup original go.mod if it exists
if exist go.mod (
    if not exist go.mod.backup (
        echo Backing up go.mod...
        copy go.mod go.mod.backup >nul
    )
)

REM Use CLI go.mod
echo Switching to CLI dependencies...
copy go_cli.mod go.mod >nul

REM Install dependencies
echo Installing dependencies...
go mod tidy

echo.
echo Building for multiple platforms...
echo.

REM Windows AMD64
echo [1/6] Building Windows AMD64...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="-s -w -X main.version=%VERSION%" -o build\windows-amd64\spotiflac.exe main_cli.go
if errorlevel 1 (
    echo ERROR: Windows AMD64 build failed
    goto :cleanup
)
echo   ✓ build\windows-amd64\spotiflac.exe

REM Windows ARM64
echo [2/6] Building Windows ARM64...
set GOOS=windows
set GOARCH=arm64
go build -ldflags="-s -w -X main.version=%VERSION%" -o build\windows-arm64\spotiflac.exe main_cli.go
if errorlevel 1 (
    echo ERROR: Windows ARM64 build failed
    goto :cleanup
)
echo   ✓ build\windows-arm64\spotiflac.exe

REM Linux AMD64
echo [3/6] Building Linux AMD64...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="-s -w -X main.version=%VERSION%" -o build\linux-amd64\spotiflac main_cli.go
if errorlevel 1 (
    echo ERROR: Linux AMD64 build failed
    goto :cleanup
)
echo   ✓ build\linux-amd64\spotiflac

REM Linux ARM64
echo [4/6] Building Linux ARM64...
set GOOS=linux
set GOARCH=arm64
go build -ldflags="-s -w -X main.version=%VERSION%" -o build\linux-arm64\spotiflac main_cli.go
if errorlevel 1 (
    echo ERROR: Linux ARM64 build failed
    goto :cleanup
)
echo   ✓ build\linux-arm64\spotiflac

REM macOS AMD64 (Intel)
echo [5/6] Building macOS AMD64 (Intel)...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="-s -w -X main.version=%VERSION%" -o build\darwin-amd64\spotiflac main_cli.go
if errorlevel 1 (
    echo ERROR: macOS AMD64 build failed
    goto :cleanup
)
echo   ✓ build\darwin-amd64\spotiflac

REM macOS ARM64 (Apple Silicon)
echo [6/6] Building macOS ARM64 (Apple Silicon)...
set GOOS=darwin
set GOARCH=arm64
go build -ldflags="-s -w -X main.version=%VERSION%" -o build\darwin-arm64\spotiflac main_cli.go
if errorlevel 1 (
    echo ERROR: macOS ARM64 build failed
    goto :cleanup
)
echo   ✓ build\darwin-arm64\spotiflac

echo.
echo Copying documentation to build directories...
for /d %%d in (build\*) do (
    copy docs\README_CLI.md "%%d\README.md" >nul 2>&1
    copy LICENSE "%%d\LICENSE" >nul 2>&1
)

echo.
echo =========================================
echo    Build Summary
echo =========================================
echo.
echo Platforms built:
echo   • Windows AMD64 (x64)
echo   • Windows ARM64 (ARM)
echo   • Linux AMD64 (x64)
echo   • Linux ARM64 (ARM)
echo   • macOS AMD64 (Intel)
echo   • macOS ARM64 (Apple Silicon)
echo.
echo Output directory: build\
echo.
dir build /ad /b
echo.

:cleanup
REM Restore original go.mod
echo Restoring original go.mod...
if exist go.mod.backup (
    move /y go.mod.backup go.mod >nul
)

echo.
echo =========================================
echo    Build Complete!
echo =========================================
echo.
echo To test: .\build\windows-amd64\spotiflac.exe --help
echo.

pause
