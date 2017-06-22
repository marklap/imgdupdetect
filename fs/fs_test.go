package fs

import (
	"os"
	"testing"

	"github.com/marklap/imgdupdetect/img"
)

func TestPath(t *testing.T) {
	tstFile := "findtest.tmp"
	_, err := os.Create(tstFile)
	if err != nil {
		t.Error(err)
	}
	defer os.Remove(tstFile)

	_, err = NewPath("lkjsdlfjalksdjflkjsadf", []Matcher{img.GIFMatch})
	if err == nil {
		t.Errorf("nonsense file was found - want: error, got: nil")
	}

	tstMatch := img.NewImageMatch([]string{tstFile})

	path, err := NewPath(".", []Matcher{tstMatch})
	if err != nil {
		t.Error(err)
	}

	paths, err := path.Find()
	if err != nil {
		t.Error(err)
	}

	if len(paths) < 1 {
		t.Errorf("incorrect number of paths found - want: 1, got: %d", len(paths))
	}

	if paths[0] != tstFile {
		t.Errorf("path mismatch - want: %s, got: %s", tstFile, paths[0])
	}

}
