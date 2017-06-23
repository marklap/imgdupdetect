package datastore

import (
	"fmt"

	"github.com/boltdb/bolt"
)

// Config is the datastore config
type Config struct {
	Path string
}

// Datastorer stores image fingerprints and their associated context data like filename and path.
type Datastorer interface {
	// Open opens the datastore for reading and writing.
	Open(Config) (*Datastore, error)

	// Close closes the datastore; no further transactions will be completed.
	Close() error

	// Get gets the set of fileData that has the same fingerprint.
	Get(collection string, fingerprint []byte) ([][]byte, error)

	// Add adds fileData for a fingerprint.
	Add(collection string, fingerprint []byte, fileData []byte) error

	// Remove removes fileData for a fingerprint.
	Remove(collection string, fingerprint []byte, fileData []byte) error
}

// Datastore is the default implementation of a Datastorer.
type Datastore struct {
	Cfg Config
	db  *bolt.DB
}

// Open opens the default datastore and preps it for transactions.
func Open(cfg Config) (*Datastore, error) {
	db, err := bolt.Open(cfg.Path, 0600, nil)
	if err != nil {
		return nil, err
	}

	return &Datastore{
		Cfg: cfg,
		db:  db,
	}, nil
}

// Close closes the default datastore.
func (d *Datastore) Close() error {
	return d.db.Close()
}

// Get gets the file data associated with this fingerprint.
func (d *Datastore) Get(col string, fp []byte) ([][]byte, error) {
	return nil, fmt.Errorf("not yet implemented")
}

// Add adds file data to the set of file data associated with this fingerprint.
func (d *Datastore) Add(col string, fp []byte, fileData []byte) error {
	return fmt.Errorf("not yet implemented")
}

// Remove removes a particular file from the set of files associated with this fingerprint.
func (d *Datastore) Remove(col string, fp []byte, data []byte) error {
	return fmt.Errorf("not yet implemented")
}
