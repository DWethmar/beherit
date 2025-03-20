package sprite

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"io/fs"
)

func LoadPng(fs fs.ReadFileFS, name string) (image.Image, error) {
	content, err := fs.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to load png: %w", err)
	}

	img, err := png.Decode(bytes.NewReader(content))
	if err != nil {
		return nil, fmt.Errorf("failed to decode png: %w", err)
	}

	return img, nil
}
