package tray

import (
	"bytes"
	"image"
	"image/png"
	"log"

	"golang.org/x/image/draw"
	"github.com/JB-SelfCompany/Tyr-Desktop/internal/resources"
)

// ResizeIcon resizes the icon to 32x32 for system tray
// System tray icons should be small (typically 16x16 or 32x32 pixels)
func ResizeIcon() []byte {
	// Decode the original PNG image
	img, err := png.Decode(bytes.NewReader(resources.ResourceTyrPng.StaticContent))
	if err != nil {
		log.Printf("Failed to decode tray icon: %v", err)
		return resources.ResourceTyrPng.StaticContent
	}

	// Create a new 32x32 image
	dst := image.NewRGBA(image.Rect(0, 0, 32, 32))

	// Resize using high-quality Catmull-Rom interpolation
	draw.CatmullRom.Scale(dst, dst.Bounds(), img, img.Bounds(), draw.Over, nil)

	// Encode back to PNG
	var buf bytes.Buffer
	if err := png.Encode(&buf, dst); err != nil {
		log.Printf("Failed to encode tray icon: %v", err)
		return resources.ResourceTyrPng.StaticContent
	}

	return buf.Bytes()
}
