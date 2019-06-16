package datastore

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	tableName = "images"
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

	// Get gets md5 for a file.
	Get(filePath string) (string, error)

	// Put adds an md5 for a file path.
	Put(filePath, hash string) error

	// // Delete removes filepath from the database.
	// Delete(filePath string) error

	// // Dump streams all objects in the database to the given file writer
	// Dump(writer io.Writer) error
}

// Datastore is the default implementation of a Datastorer.
type Datastore struct {
	Cfg Config
	db  *sql.DB
}

// Open opens the default datastore and preps it for transactions.
func Open(cfg Config) (*Datastore, error) {
	db, err := sql.Open("sqlite3", cfg.Path)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(
		"CREATE TABLE IF NOT EXISTS images (filePath TEXT NOT NULL PRIMARY KEY, md5sum TEXT NOT NULL)",
	)
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

// Get gets the md5sum for a file
func (d *Datastore) Get(filePath string) (string, error) {
	rows, err := d.db.Query("SELECT md5sum FROM images WHERE filePath = ?", filePath)
	if err != nil {
		return "", nil
	}

	if rows.Next() {
		var md5sum string
		if err = rows.Scan(&md5sum); err != nil {
			return md5sum, nil
		}
	}

	return "", nil
}

// Put puts a filepath and md5sum in the database
func (d *Datastore) Put(filePath, md5sum string) error {
	_, err := d.db.Exec("INSERT OR REPLACE INTO images (filePath, md5sum) VALUES (?, ?)", filePath, md5sum)
	if err != nil {
		return err
	}
	return nil
}
