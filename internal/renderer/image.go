package renderer

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"strings"

	"github.com/disintegration/imaging"
	_ "golang.org/x/image/webp"
)

// ImageRenderer renders images to terminal using Unicode half-blocks
type ImageRenderer struct {
	maxWidth  int
	maxHeight int
}

// NewImageRenderer creates a new image renderer
func NewImageRenderer(maxWidth, maxHeight int) *ImageRenderer {
	return &ImageRenderer{
		maxWidth:  maxWidth,
		maxHeight: maxHeight,
	}
}

// RenderImage renders an image as Unicode blocks
func (r *ImageRenderer) RenderImage(imageData []byte) (string, error) {
	// Decode image
	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		return "", fmt.Errorf("failed to decode image: %w", err)
	}

	// Calculate dimensions (each character represents 2 vertical pixels using half-blocks)
	bounds := img.Bounds()
	imgWidth := bounds.Dx()
	imgHeight := bounds.Dy()

	// Calculate target dimensions to fit within terminal
	// Each character cell is roughly 2 pixels tall when using half-blocks
	targetWidth := r.maxWidth
	targetHeight := r.maxHeight * 2

	// Maintain aspect ratio
	if imgWidth > targetWidth || imgHeight > targetHeight {
		ratio := float64(imgWidth) / float64(imgHeight)
		if imgWidth > targetWidth {
			targetWidth = r.maxWidth
			targetHeight = int(float64(targetWidth) / ratio)
		}
		if targetHeight > r.maxHeight*2 {
			targetHeight = r.maxHeight * 2
			targetWidth = int(float64(targetHeight) * ratio)
		}
	} else {
		targetWidth = imgWidth
		targetHeight = imgHeight
	}

	// Ensure even height for half-block rendering
	if targetHeight%2 != 0 {
		targetHeight++
	}

	// Resize image
	resized := imaging.Resize(img, targetWidth, targetHeight, imaging.Lanczos)

	// Render using half-blocks (▀ for upper half)
	var out strings.Builder

	// Add image info
	out.WriteString(fmt.Sprintf("Image: %dx%d (displayed as %dx%d)\n\n", imgWidth, imgHeight, targetWidth, targetHeight/2))

	// Process pairs of rows
	for y := 0; y < targetHeight; y += 2 {
		for x := 0; x < targetWidth; x++ {
			// Get colors for upper and lower pixels
			upperColor := resized.At(x, y)
			var lowerColor color.Color
			if y+1 < targetHeight {
				lowerColor = resized.At(x, y+1)
			} else {
				lowerColor = color.RGBA{0, 0, 0, 0}
			}

			// Convert to RGB
			ur, ug, ub, ua := upperColor.RGBA()
			lr, lg, lb, la := lowerColor.RGBA()

			// Convert from uint32 (0-65535) to uint8 (0-255)
			upperR, upperG, upperB := uint8(ur>>8), uint8(ug>>8), uint8(ub>>8)
			lowerR, lowerG, lowerB := uint8(lr>>8), uint8(lg>>8), uint8(lb>>8)
			upperA := uint8(ua >> 8)
			lowerA := uint8(la >> 8)

			// Handle transparency
			if upperA < 128 && lowerA < 128 {
				out.WriteString(" ")
			} else if upperA < 128 {
				// Only lower pixel is visible - use full block with lower color
				out.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm█\x1b[0m", lowerR, lowerG, lowerB))
			} else if lowerA < 128 {
				// Only upper pixel is visible - use upper half block with upper color
				out.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%dm▀\x1b[0m", upperR, upperG, upperB))
			} else {
				// Both pixels visible - use half block with upper as foreground, lower as background
				out.WriteString(fmt.Sprintf("\x1b[38;2;%d;%d;%d;48;2;%d;%d;%dm▀\x1b[0m",
					upperR, upperG, upperB, lowerR, lowerG, lowerB))
			}
		}
		out.WriteString("\n")
	}

	return out.String(), nil
}

// IsImageMIME checks if a MIME type is an image
func IsImageMIME(mimeType string) bool {
	return strings.HasPrefix(mimeType, "image/png") ||
		strings.HasPrefix(mimeType, "image/jpeg") ||
		strings.HasPrefix(mimeType, "image/jpg") ||
		strings.HasPrefix(mimeType, "image/gif") ||
		strings.HasPrefix(mimeType, "image/webp") ||
		strings.HasPrefix(mimeType, "image/bmp")
}
