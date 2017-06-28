package datastore

import (
	"bytes"
	"testing"
)

var (
	tstDatastorePath = "./testdata.dstore"
	tstCollection    = "test"
)

func TestDataStore(t *testing.T) {
	ds, err := Open(Config{tstDatastorePath})
	if err != nil {
		t.Error(err)
	}

	tstHash := []byte{1, 2, 3}
	tstFileName := "/tmp/my/file/name.jpg"
	wantKey := "key"
	wantValue := []byte("value")
	err = ds.Add(tstCollection, tstHash, tstFileName, map[string][]byte{wantKey: wantValue})
	if err != nil {
		t.Error(err)
	}

	col, err := ds.Get(tstCollection, tstHash)
	if err != nil {
		t.Error(err)
	}

	if f, found := col[tstFileName]; !found {
		t.Errorf("filename %s does not exist in collection", tstFileName)
		if gotValue, found := f[wantKey]; !found {
			t.Errorf("key %s not found", wantKey)
			if bytes.Equal(wantValue, gotValue) {
				t.Errorf("value mismatch - want: %s, got: %s", wantValue, gotValue)
			}
		}
	}

	err = ds.Remove(tstCollection, tstHash, tstFileName)
	if err != nil {
		t.Error(err)
	}
}
