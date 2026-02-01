@echo off

REM Change to project root
cd /d "%~dp0\.."

echo =========================================
echo    SpotiFLAC CLI Build Script
echo =========================================
echo.

REM Clean old binary
if exist spotiflac.exe del spotiflac.exe

REM Backup original go.mod
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

REM Build for Windows
echo.
echo Building for Windows AMD64...
go build -ldflags="-s -w" -o spotiflac.exe main_cli.go

if errorlevel 1 (
    echo.
    echo ERROR: Build failed
    goto :restore
)

REM Restore original go.mod
:restore
echo.
echo Restoring original go.mod...
if exist go.mod.backup (
    move /y go.mod.backup go.mod >nul
)

if not errorlevel 1 (
    echo.
    echo =========================================
    echo    Build Complete!
    echo =========================================
    echo.
    echo Binary: spotiflac.exe
    echo Run: spotiflac.exe --help
)

echo.
pause
