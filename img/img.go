package img

import (
	"crypto/sha256"
	"encoding/binary"
	"image"
	"math"
	"os"

	_ "image/gif"  // registers gif encoding
	_ "image/jpeg" // registers jpg encoding
	_ "image/png"  // registers png encoding

	log "github.com/sirupsen/logrus"
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
	Path     string
	Type     string
	Config   image.Config
	FileInfo os.FileInfo
}

// NewImage creates a new Image
func NewImage(path string) (*Image, error) {
	log.Debugf("creating Image for path: %s", path)

	fd, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	fi, err := fd.Stat()
	if err != nil {
		return nil, err
	}

	imgCfg, imgType, err := image.DecodeConfig(fd)
	if err != nil {
		return nil, err
	}

	return &Image{
		Path:     path,
		Type:     imgType,
		Config:   imgCfg,
		FileInfo: fi,
	}, nil
}

// midPoints find the middle
func midPoints(w, h int) (x, y int) {
	return int(math.Floor(float64(w) / 2.0)), int(math.Floor(float64(h) / 2.0))
}

// FingerPrint returns a unique fingerprint for an image - well... for practical purposes
func (i *Image) FingerPrint() ([]byte, error) {
	buf := make([]byte, (i.Config.Width+i.Config.Height)*8) // 8 bytes for size + (2 bytes per color (0xffff), 4 colors in a pixel (rgba), 8 bytes per pixel)

	fd, err := os.Open(i.Path)
	if err != nil {
		return nil, err
	}
	defer fd.Close()

	image, _, err := image.Decode(fd)
	if err != nil {
		return nil, err
	}

	bounds := image.Bounds()
	midX, midY := midPoints(i.Config.Width, i.Config.Height)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		startOffset := (y - bounds.Min.Y) * 8 // the start of the slice

		r, g, b, a := image.At(midX, y).RGBA()
		for j, c := range []uint32{r, g, b, a} {
			offset := startOffset + (j * 2)
			buf[offset] = byte(c >> 8)
			buf[offset+1] = byte(c)
		}
	}

	for x := bounds.Min.X; x < bounds.Max.X; x++ {
		startOffset := (i.Config.Height + x - bounds.Min.X) * 8 // the start of the slice

		r, g, b, a := image.At(x, midY).RGBA()
		for j, c := range []uint32{r, g, b, a} {
			offset := startOffset + (j * 2)
			buf[offset] = byte(c >> 8)
			buf[offset+1] = byte(c)
		}
	}

	sum := sha256.Sum256(buf)
	res := make([]byte, len(sum))
	copy(res, sum[:])
	return res, nil
}

// Size returns the size
func (i *Image) Size() uint64 {
	return uint64(i.FileInfo.Size())
}

// Height returns the height of the image
func (i *Image) Height() uint64 {
	return uint64(i.Config.Height)
}

// Width returns the width of the image
func (i *Image) Width() uint64 {
	return uint64(i.Config.Width)
}

// SizeByteSlice returns the size of the image in a byte array
func (i *Image) SizeByteSlice() []byte {
	return uint64ToByteSlice(i.Size())
}

// HeightByteSlice returns the size of the image in a byte array
func (i *Image) HeightByteSlice() []byte {
	return uint64ToByteSlice(i.Height())
}

// WidthByteSlice returns the size of the image in a byte array
func (i *Image) WidthByteSlice() []byte {
	return uint64ToByteSlice(i.Width())
}

// uint64ToByteSlice returns a byte slice representing the int
func uint64ToByteSlice(i uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

// FingerPrintCollection is a collection of fingerprints
type FingerPrintCollection struct {
	FingerPrints []FingerPrint
}

// FingerPrint is a fingerprint and the files associated with it
type FingerPrint struct {
	Hash   []byte
	Images []string
}
