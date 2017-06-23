package img

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	tstImagePath   = filepath.Clean(filepath.Join(filepath.Dir(os.Args[0]), "..", "..", "..", ".."))
	tstImageOrig   = filepath.Join(tstImagePath, "monkey.orig.jpg")
	tstImageCopy   = filepath.Join(tstImagePath, "monkey.dup.jpg")
	tstImageGrow   = filepath.Join(tstImagePath, "monkey.lg.jpg")
	tstImageShrink = filepath.Join(tstImagePath, "monkey.sm.jpg")
	tstImageCrop   = filepath.Join(tstImagePath, "monkey.crop.jpg")
	tstImageSharp  = filepath.Join(tstImagePath, "monkey.sharp10.jpg")
	tstImageWidth  = 1600
	tstImageHeight = 1200
)

func TestImageMatch(t *testing.T) {
	want := []string{"findtest.tmp"}
	got := NewImageMatch(want).Patterns()
	if len(want) != len(got) && want[0] != want[1] {
		t.Errorf("mismatched patterns - want: %s, got: %s", want, got)
	}
}

func TestNewImage(t *testing.T) {
	wantT := "jpeg"
	wantW := tstImageWidth
	wantH := tstImageHeight

	got, err := NewImage(tstImageOrig)
	if err != nil {
		t.Error(err)
	}

	if wantT != got.Type {
		t.Errorf("image type mismatch - want: %s, got: %s", wantT, got.Type)
	}

	if wantW != got.Config.Width {
		t.Errorf("width mismatch - want: %d, got: %d", wantW, got.Config.Width)
	}

	if wantH != got.Config.Height {
		t.Errorf("height mismatch - want: %d, got: %d", wantH, got.Config.Height)
	}
}

func TestFingerPrint(t *testing.T) {
	wantIface := reflect.TypeOf((*FingerPrinter)(nil)).Elem()
	gotIface := reflect.TypeOf((*Image)(nil))
	if !gotIface.Implements(wantIface) {
		t.Errorf("does not implement interface - want: %s, got: %s", wantIface, gotIface)
	}

	orig, err := NewImage(tstImageOrig)
	if err != nil {
		t.Error(err)
	}

	origFp, err := orig.FingerPrint()
	if err != nil {
		t.Error(err)
	}

	origSum := fmt.Sprintf("%x", origFp)

	for _, tstImg := range []struct {
		path string
		want bool
	}{
		{tstImageCopy, true},
		{tstImageCrop, false},
		{tstImageGrow, false},
		{tstImageSharp, false},
		{tstImageShrink, false},
	} {
		i, err := NewImage(tstImg.path)
		if err != nil {
			t.Error(err)
		}

		f, err := i.FingerPrint()
		if err != nil {
			t.Error(err)
		}

		s := fmt.Sprintf("%x", f)
		if got := origSum == s; got != tstImg.want {
			t.Errorf("duplicate detection failed - want: %t, got: %t", tstImg.want, got)
		}
	}
}

func BenchmarkFingerPrint(b *testing.B) {
	var img *Image
	var err error
	img, err = NewImage(tstImageOrig)
	if err != nil {
		b.Error(err)
	}
	for i := 0; i < b.N; i++ {
		_, err = img.FingerPrint()
		if err != nil {
			b.Error(err)
		}
	}
}
