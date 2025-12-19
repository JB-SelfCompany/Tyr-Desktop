package main

// version is the application version
// This can be set at build time using ldflags:
//   go build -ldflags "-X main.version=x.y.z"
// For Wails builds, this will be automatically set from wails.json
var version = "dev"