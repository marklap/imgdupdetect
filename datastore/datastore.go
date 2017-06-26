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
	Get(collection string, fingerprint []byte) (map[string][]map[string][]byte, error)

	// Add adds a file data for a fingerprint.
	Add(collection string, fingerprint []byte, filename string, data map[string][]byte) error

	// Remove removes a file from the set of files for this fingerprint.
	Remove(collection string, fingerprint []byte, filename string) error
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

// Get gets the set of file data associated with this fingerprint.
func (d *Datastore) Get(col string, fp []byte) (map[string][]map[string][]byte, error) {
	// var res map[string][]map[string][]byte
	// var err error
	// err = d.db.View(func(tx *bolt.Tx) error {
	// 	root := tx.Bucket([]byte(col))
	// 	if root == nil {
	// 		return nil
	// 	}
	// 	bkt := root.Bucket(fp)
	// 	if bkt == nil {
	// 		return nil
	// 	}
	// 	res = make(map[string][]map[string][]byte)
	// 	bkt.ForEach(func(k, v []byte) error {
	// 		var m = make(map[string][]byte)
	// 		b.ForEach(func(k, v []byte) error {
	// 			m[string(k)] = v
	// 		})
	// 		res[string(n)] = append(res[string(n)], m)
	// 	})
	// })
	// return res, err
	return nil, fmt.Errorf("not yet implemented")
}

// Add adds file data to the set of file data associated with this fingerprint.
func (d *Datastore) Add(col string, fp []byte, name string, data map[string][]byte) error {
	var err error
	err = d.db.Update(func(tx *bolt.Tx) error {
		cBkt, err := tx.CreateBucketIfNotExists([]byte(col))
		if err != nil {
			return err
		}

		fpBkt, err := cBkt.CreateBucketIfNotExists(fp)
		if err != nil {
			return err
		}

		fileBkt, err := fpBkt.CreateBucketIfNotExists([]byte(name))
		if err != nil {
			return err
		}

		for k, v := range data {
			fileBkt.Put([]byte(k), v)
		}

		return nil
	})
	return err
}

// Remove removes a particular file from the set of files associated with this fingerprint.
func (d *Datastore) Remove(col string, fp []byte, name string) error {
	return fmt.Errorf("not yet implemented")
}
