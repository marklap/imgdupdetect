package datastore

import (
	"fmt"
	"io"

	"github.com/boltdb/bolt"
	log "github.com/sirupsen/logrus"
)

var (
	errTmplBucketNotFound = "bucket not found: %s"
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
	Get(collection string, fingerprint []byte) (map[string]map[string][]byte, error)

	// Add adds a file data for a fingerprint.
	Add(collection string, fingerprint []byte, filename string, data map[string][]byte) error

	// Remove removes a file from the set of files for this fingerprint.
	Remove(collection string, fingerprint []byte, filename string) error

	// Dump streams all objects in the database to the given file writer
	Dump(writer io.Writer) error
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
func (d *Datastore) Get(col string, fp []byte) (map[string]map[string][]byte, error) {
	var res = make(map[string]map[string][]byte)
	err := d.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(col))
		if root == nil {
			return fmt.Errorf(errTmplBucketNotFound, col)
		}

		fpBkt := root.Bucket(fp)
		if fpBkt == nil {
			return fmt.Errorf(errTmplBucketNotFound, fp)
		}

		fpBkt.ForEach(func(k, v []byte) error {
			if v == nil {
				res[string(k)] = make(map[string][]byte)
				fileBkt := fpBkt.Bucket(k)
				if fileBkt == nil {
					log.Errorf(errTmplBucketNotFound, k)
				}
				fileBkt.ForEach(func(mKey, mVal []byte) error {
					res[string(k)][string(mKey)] = mVal
					return nil
				})
			}
			return nil
		})
		return nil
	})
	return res, err
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
	err := d.db.Update(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(col))
		if root == nil {
			return fmt.Errorf(errTmplBucketNotFound, col)
		}

		fpBkt := root.Bucket(fp)
		if fpBkt == nil {
			return fmt.Errorf(errTmplBucketNotFound, fp)
		}

		err := fpBkt.DeleteBucket([]byte(name))
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

// GetFingerPrints gets the fingerprints
func (d *Datastore) GetFingerPrints(col string) [][]byte {
	var res [][]byte
	d.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(col))
		root.ForEach(func(k, v []byte) error {
			if v == nil {
				res = append(res, k)
			}
			return nil
		})
		return nil
	})
	return res
}

// GetImages gets the images associated with a fingerprint
func (d *Datastore) GetImages(col string, fp []byte) []string {
	var res []string
	d.db.View(func(tx *bolt.Tx) error {
		root := tx.Bucket([]byte(col))

		fpBkt := root.Bucket(fp)
		fpBkt.ForEach(func(k, v []byte) error {
			if v == nil {
				res = append(res, string(k))
			}
			return nil
		})
		return nil
	})
	return res
}

// Dump dumps the contents of the database to the provided file writer
func (d *Datastore) Dump(w io.Writer) error {
	return d.db.View(func(tx *bolt.Tx) error {
		return tx.ForEach(func(name []byte, _ *bolt.Bucket) error {
			b := tx.Bucket(name)
			c := b.Cursor()
			for k, v := c.First(); k != nil; k, v = c.Next() {
				w.Write([]byte(fmt.Sprintf("bucket=%s, key=%s, value=%s\n", string(name), string(k), string(v))))
			}
			return nil
		})
	})
}
