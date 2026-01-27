// Package version provides centralized version information for the application.
// The Version variable can be overridden at build time using ldflags:
//
//	go build -ldflags "-X github.com/JB-SelfCompany/Tyr-Desktop/internal/version.Version=x.x.x"
//
// For Wails builds:
//
//	wails build -ldflags "-X github.com/JB-SelfCompany/Tyr-Desktop/internal/version.Version=x.x.x"
package version

// Version is the application version.
// This value is set at compile time via ldflags or manually updated here.
var Version = "2.1.1"
