package fs

import (
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

var (
	errRootNotDir = fmt.Errorf("path is not a directory")
)

// Imager specifies an image
type Imager interface {
	// Name returns the name of the Image type
	Name() string
	// Patterns returns the patterns associated with this Image type
	Patterns() []string
}

// Image helps with finding images of the specified name.
type Image struct {
	name     string
	patterns []string
}

// NewImage creates a new Image
func NewImage(n string, p []string) *Image {
	return &Image{
		name:     n,
		patterns: p,
	}
}

// Name returns the name of this Img
func (i *Image) Name() string {
	return i.name
}

// Patterns returns the patterns associated with this image
func (i *Image) Patterns() []string {
	return i.patterns
}

var (
	// GIF is a gif image
	GIF = &Image{name: "gif", patterns: []string{"*.gif"}}
	// JPG is a jpg image
	JPG = &Image{name: "jpg", patterns: []string{"*.jpg", "*.jpeg"}}
	// PNG is a png image
	PNG = &Image{name: "png", patterns: []string{"*.png"}}
)

// Path is used to search for images in the specified root path.
type Path struct {
	Name         string
	Root         *os.File
	RootFileInfo os.FileInfo
	Recursive    bool
}

// NewPath creates a new path
func NewPath(root string, recurs bool) (*Path, error) {
	fi, err := os.Stat(root)
	if err != nil {
		return nil, err
	}

	if !fi.IsDir() {
		return nil, errRootNotDir
	}

	fd, err := os.Open(root)
	if err != nil {
		return nil, err
	}

	return &Path{
		Name:         root,
		Root:         fd,
		RootFileInfo: fi,
		Recursive:    recurs,
	}, nil
}

// FindImages finds all imgs in the specified dir directory and returns a list of file paths.
func (p *Path) FindImages(imgs []Imager) ([]string, error) {
	var paths = []string{}
	err := filepath.Walk(p.Root.Name(), func(path string, info os.FileInfo, err error) error {
		for _, img := range imgs {
			for _, pattern := range img.Patterns() {
				if matched, err := filepath.Match(pattern, filepath.Base(path)); err == nil {
					if matched {
						paths = append(paths, path)
						log.Debugf("path matched on pattern %s: %s", pattern, path)
					} else {
						log.Debugf("path NO MATCH on pattern %s: %s", pattern, path)
					}
				} else {
					log.Error(err)
				}
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return paths, nil
}
