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

// Matcher specifies an method of finding files
type Matcher interface {
	// Patterns returns the patterns associated with this Image type
	Patterns() []string
}

// Path is used to search for images in the specified root path.
type Path struct {
	Name         string
	Root         *os.File
	RootFileInfo os.FileInfo
	Matchers     []Matcher
}

// NewPath creates a new path
func NewPath(root string, matchers []Matcher) (*Path, error) {
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
		Matchers:     matchers,
	}, nil
}

// Find finds all files in the specified dir directory and returns a list of paths.
func (p *Path) Find() ([]string, error) {
	var paths = []string{}
	err := filepath.Walk(p.Root.Name(), func(path string, info os.FileInfo, err error) error {
		for _, match := range p.Matchers {
			for _, pattern := range match.Patterns() {
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
