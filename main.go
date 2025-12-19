package main

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/windows"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Check if app should start minimized
	startMinimized := false
	for _, arg := range os.Args[1:] {
		if arg == "--minimized" {
			startMinimized = true
			log.Println("Starting in minimized mode (from autostart)")
			break
		}
	}

	// Create application instance
	app := NewApp()
	app.startMinimized = startMinimized

	// Get executable directory for WebView2 user data path
	// WebView2 will automatically create "EBWebView" folder inside this path
	exePath, err := os.Executable()
	if err != nil {
		log.Printf("Warning: Failed to get executable path: %v", err)
		exePath = ""
	}
	webviewDataPath := filepath.Dir(exePath)

	// Run Wails application with Wails v2 configuration
	err = wails.Run(&options.App{
		// Window configuration
		Title:     "Tyr Desktop - P2P Email Client",
		Width:     1200,
		Height:    800,
		MinWidth:  1080, // Increased to prevent layout issues on small screens
		MinHeight: 700,  // Increased to ensure content fits properly

		// StartHidden: When true, application is hidden until WindowShow is called
		// This is used for autostart - app starts in system tray without showing window
		StartHidden: startMinimized,

		// Asset server configuration
		AssetServer: &assetserver.Options{
			Assets: assets,
		},

		// Background color (Y2K Dark Theme background)
		BackgroundColour: &options.RGBA{R: 10, G: 14, B: 26, A: 1},

		// Wails v2 Lifecycle hooks (lowercase method names)
		OnStartup:     app.startup,
		OnDomReady:    app.domReady,
		OnBeforeClose: app.beforeClose,
		OnShutdown:    app.shutdown,

		// Bind the App struct to the frontend
		// This makes all exported methods available in TypeScript
		Bind: []interface{}{
			app,
		},

		// Windows-specific options
		Windows: &windows.Options{
			WebviewIsTransparent: false,
			WindowIsTranslucent:  false,
			DisableWindowIcon:    false,
			WebviewUserDataPath:  webviewDataPath, // Store WebView2 data in installation directory instead of AppData
			// ContentProtection:    true, // Uncomment to block screen capture (privacy)
		},

		// Uncomment to enable devtools in production
		// Debug: options.Debug{
		// 	OpenInspectorOnStartup: false,
		// },
	})

	if err != nil {
		log.Fatalf("Error running application: %v", err)
	}
}
