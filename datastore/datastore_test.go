package datastore

import (
	"testing"
)

var (
	tstDatastorePath = "./testdata.dstore"
)

func TestDataStoreAdd(t *testing.T) {
	ds, err := Open(Config{tstDatastorePath})
	if err != nil {
		t.Error(err)
	}

	err = ds.Add("test", []byte{1, 2, 3}, "/tmp/my/file/name.jpg", map[string][]byte{"key": []byte("value")})
	if err != nil {
		t.Error(err)
	}
}
