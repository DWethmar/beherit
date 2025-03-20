package sprite_test

import (
	"bytes"
	"image"
	"image/color"
	"image/png"
	"testing"
	"testing/fstest"

	"github.com/dwethmar/beherit/cmd/mygame/sprite"
)

func createTestPNG() ([]byte, error) {
	img := image.NewRGBA(image.Rect(0, 0, 10, 10))
	for x := range 10 {
		for y := range 10 {
			img.Set(x, y, color.RGBA{uint8(x * 25), uint8(y * 25), 0, 255})
		}
	}
	var buf bytes.Buffer
	err := png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Test_loadPng(t *testing.T) {
	// Create a test PNG image
	pngData, err := createTestPNG()
	if err != nil {
		t.Fatalf("Failed to create test PNG image: %v", err)
	}

	// Create an in-memory filesystem with the PNG data
	testFS := fstest.MapFS{
		"test.png": &fstest.MapFile{
			Data: pngData,
		},
	}

	// Call loadPng with the in-memory filesystem
	img, err := sprite.LoadPng(testFS, "test.png")
	if err != nil {
		t.Fatalf("loadPng failed: %v", err)
	}
	if img == nil {
		t.Fatal("loadPng returned nil image")
	}

	// Verify image dimensions
	bounds := img.Bounds()
	if bounds.Dx() != 10 || bounds.Dy() != 10 {
		t.Errorf("Expected image dimensions 10x10, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}
