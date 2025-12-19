@echo off
REM Tyr Desktop - Windows Build Script
REM This script builds the application for Windows
REM Requirements: Go 1.21+ installed and in PATH

echo ========================================
echo Tyr Desktop - Windows Build Script
echo ========================================
echo.

REM Check if Go is installed
where go >nul 2>nul
if %errorlevel% neq 0 (
    echo ERROR: Go is not installed or not in PATH!
    echo Please install Go from https://golang.org/dl/
    echo Minimum required version: 1.21
    exit /b 1
)

REM Display Go version
echo Checking Go installation...
go version
echo.

REM Check if Node.js/npm is installed
where node >nul 2>nul
if %errorlevel% neq 0 (
    echo ERROR: Node.js is not installed or not in PATH!
    echo Please install Node.js from https://nodejs.org/
    echo This will also install npm (Node Package Manager)
    exit /b 1
)

REM Display Node.js and npm versions
echo Checking Node.js installation...
node --version
npm --version
echo.

REM Install/update Wails CLI if not installed
echo Checking Wails CLI installation...
where wails >nul 2>nul
if %errorlevel% neq 0 (
    echo Wails CLI is not installed. Installing now...
    echo This may take a few minutes...
    go install github.com/wailsapp/wails/v2/cmd/wails@latest
    if %errorlevel% neq 0 (
        echo ERROR: Failed to install Wails CLI
        echo Please check your Go installation and try again
        exit /b 1
    )
    echo Wails CLI installed successfully!
    echo.
) else (
    echo Wails CLI is already installed
    wails version
    echo.
)

REM Download Go module dependencies
echo Downloading Go module dependencies...
go mod download
if %errorlevel% neq 0 (
    echo ERROR: Failed to download Go dependencies
    exit /b 1
)
echo Go dependencies downloaded successfully!
echo.

REM Extract version from wails.json
echo Reading version from wails.json...
for /f "tokens=2 delims=:, " %%a in ('findstr /r "\"version\":" wails.json') do (
    set VERSION=%%a
)
set VERSION=%VERSION:"=%
echo Building version: %VERSION%
echo.

echo Step 1: Preparing build directory...
REM Clean build directories
if exist build\bin rmdir /s /q build\bin
if exist build\windows rmdir /s /q build\windows
if exist frontend\dist rmdir /s /q frontend\dist
echo Cleaned previous builds.

REM Create build directory structure
if not exist build mkdir build
if not exist build\windows mkdir build\windows

REM Copy build assets from assets/build
echo Copying build assets from assets/build...
if exist assets\build\appicon.png (
    copy /Y assets\build\appicon.png build\appicon.png >nul
    echo   - Copied appicon.png
) else (
    echo   WARNING: assets\build\appicon.png not found!
)
if exist assets\build\windows\icon.ico (
    copy /Y assets\build\windows\icon.ico build\windows\icon.ico >nul
    echo   - Copied icon.ico
) else (
    echo   WARNING: assets\build\windows\icon.ico not found!
)
echo Done.
echo.

echo Step 2: Installing frontend dependencies...
cd frontend
call npm install
if %errorlevel% neq 0 (
    echo ERROR: Failed to install frontend dependencies
    exit /b 1
)
cd ..
echo Done.
echo.

echo Step 3: Building application...
echo This may take several minutes...
wails build -clean -platform windows/amd64 -webview2 download -ldflags "-X main.version=%VERSION%"
if %errorlevel% neq 0 (
    echo ERROR: Build failed
    exit /b 1
)
echo Done.
echo.

echo Step 4: Updating system tray icon...
if exist build\windows\icon.ico (
    copy /Y build\windows\icon.ico internal\resources\tyr.ico >nul
    echo System tray icon updated
)
echo Done.
echo.

echo ========================================
echo Build completed successfully!
echo ========================================
if exist build\bin\Tyr-Desktop.exe (
    echo Executable: build\bin\Tyr-Desktop.exe
) else (
    echo ERROR: Executable not found at build\bin\Tyr-Desktop.exe
)
echo.

pause
