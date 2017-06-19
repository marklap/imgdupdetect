package fs

import (
	"testing"
)

func TestImage(t *testing.T) {
	for _, i := range []struct {
		name  string
		image Imager
	}{
		{"gif", GIF},
		{"jpg", JPG},
		{"png", PNG},
	} {
		if i.name != i.image.Name() {
			t.Errorf("name mismatch - want: %s, got: %s", i.name, i.image.Name())
		}
	}
}

func TestPath(t *testing.T) {
	_, err := NewPath("lkjsdlfjalksdjflkjsadf", true)
	if err == nil {
		t.Errorf("nonsense file was found - want: error, got: nil")
	}

	_, err = NewPath(".", false)
	if err != nil {
		t.Error(err)
	}
}
