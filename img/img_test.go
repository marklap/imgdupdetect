package img

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

var (
	tstImagePath   = filepath.Clean(filepath.Join(filepath.Dir(os.Args[0]), "..", "..", "..", "..", "monkey.orig.jpg"))
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

	got, err := NewImage(tstImagePath)
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

	i, err := NewImage(tstImagePath)
	if err != nil {
		t.Error(err)
	}

	fp, err := i.FingerPrint()
	if err != nil {
		t.Error(err)
	}

	t.Logf("%x", fp)

}

func BenchmarkFingerPrint(b *testing.B) {
	var img *Image
	var err error
	img, err = NewImage(tstImagePath)
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
