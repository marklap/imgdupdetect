package img

import (
	"fmt"
	"image"
	"os"
)

// FingerPrinter specifies a value that can generate a practically-unique fingerprint for an image
type FingerPrinter interface {
	FingerPrint() ([]byte, error)
}

// ImageMatch matches typical image files
type ImageMatch struct {
	patterns []string
}

// NewImageMatch creates a new ImageMatch
func NewImageMatch(patterns []string) *ImageMatch {
	return &ImageMatch{
		patterns: patterns,
	}
}

// Patterns returns the patters for this Matcher
func (i *ImageMatch) Patterns() []string {
	return i.patterns
}

var (
	// GIFMatch matches on gif files
	GIFMatch = &ImageMatch{[]string{"*.gif"}}
	// JPGMatch matches jpegs
	JPGMatch = &ImageMatch{[]string{"*.jpg", "*.jpeg"}}
	// PNGMatch matches png files
	PNGMatch = &ImageMatch{[]string{"*.png"}}
)

// Image represents an image file
type Image struct {
	Path   string
	Type   string
	Config interface{}
}

// NewImage creates a new Image
func NewImage(path string) (*Image, error) {
	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	imgCfg, imgType, err := image.DecodeConfig(fd)
	if err != nil {
		return nil, err
	}

	return &Image{
		Path:   path,
		Type:   imgType,
		Config: imgCfg,
	}, nil
}

// FingerPrint returns a unique fingerprint for an image - well... for practical purposes
func (i *Image) FingerPrint(path string) ([]byte, error) {
	return nil, fmt.Errorf("not implemented")
}
