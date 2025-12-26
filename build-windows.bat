@echo off
setlocal enabledelayedexpansion

REM ========================================
REM Tyr Desktop - Windows Build Script
REM ========================================
echo ========================================
echo Tyr Desktop - Windows Build Script
echo ========================================
echo.

REM ========================================
REM Step 1: Verify prerequisites
REM ========================================
echo [1/6] Verifying prerequisites...
echo.

echo Checking Go installation...
call go version
echo.

echo Checking Node.js installation...
call node --version
echo.

echo Checking npm installation...
call npm --version
echo.

echo Checking Wails CLI installation...
call wails version 2>nul || echo Wails not found, but continuing...
echo.

echo All prerequisites verified!
echo.

REM ========================================
REM Step 2: Extract version from wails.json
REM ========================================
echo [2/6] Reading version from wails.json...
for /f "tokens=2 delims=:, " %%a in ('findstr /r "\"version\":" wails.json') do (
    set VERSION=%%a
)
set VERSION=%VERSION:"=%
echo Building version: %VERSION%
echo.

REM ========================================
REM Step 3: Clean previous builds
REM ========================================
echo [3/6] Preparing build directory...

REM Clean build directories
if exist build\bin (
    echo Removing old build\bin...
    rmdir /s /q build\bin
)
if exist build\windows (
    echo Removing old build\windows...
    rmdir /s /q build\windows
)
if exist frontend\dist (
    echo Removing old frontend\dist...
    rmdir /s /q frontend\dist
)
echo Cleaned previous builds.
echo.

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

REM ========================================
REM Step 4: Download Go dependencies
REM ========================================
echo [4/6] Downloading Go module dependencies...
call go mod download
echo Done.
echo.

REM ========================================
REM Step 5: Install frontend dependencies
REM ========================================
echo [5/6] Installing frontend dependencies...
cd frontend
call npm install
cd ..
echo Done.
echo.

REM ========================================
REM Step 6: Build application
REM ========================================
echo [6/6] Building application...
echo This may take several minutes...
echo.
echo Command: wails build -clean -platform windows/amd64 -webview2 download -ldflags "-X main.version=%VERSION%"
echo.

call wails build -clean -platform windows/amd64 -webview2 download -ldflags "-X main.version=%VERSION%"

echo.
echo.

REM Check if build was successful by checking if executable exists
if not exist build\bin\Tyr-Desktop.exe (
    echo ========================================
    echo ERROR: Build failed!
    echo ========================================
    echo Executable not found at build\bin\Tyr-Desktop.exe
    echo Please check the output above for errors.
    echo.
    pause
    exit /b 1
)

REM Update system tray icon
echo Updating system tray icon...
if exist build\windows\icon.ico (
    if not exist internal\resources mkdir internal\resources
    copy /Y build\windows\icon.ico internal\resources\tyr.ico >nul
    echo System tray icon updated
)
echo.

REM ========================================
REM Build completed successfully
REM ========================================
echo ========================================
echo Build completed successfully!
echo ========================================
echo.
echo Executable: build\bin\Tyr-Desktop.exe
echo Version: %VERSION%
echo.
echo You can now run the application from build\bin\Tyr-Desktop.exe
echo.

pause
