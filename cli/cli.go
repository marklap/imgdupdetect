package cli

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/marklap/imgdupdetect/datastore"
	"github.com/marklap/imgdupdetect/fs"
	"github.com/marklap/imgdupdetect/img"
	"github.com/marklap/imgdupdetect/stats"

	log "github.com/sirupsen/logrus"
)

// DupeDetectConfig is the duplicate detector CLI config
type DupeDetectConfig struct {

	// Dirs is the directories to scan for duplicates
	Dirs []string
	// Datastore is the datastore
	Datastore *datastore.Datastore
	// FingerPrintCol is the name of the collection to use for fingerprints
}

func ensureDir(path string) error {
	dirPath := filepath.Dir(path)
	if stat, err := os.Stat(dirPath); os.IsExist(err) && stat.IsDir() {
		return nil
	}

	mode := os.FileMode(0755)
	err := os.MkdirAll(dirPath, mode)
	if err != nil {
		log.Error(err)
		return err
	}
	return nil
}

func MD5Sum(cfg DupeDetectConfig) error {
	scanStats := stats.NewScanStats()

	log.Info("Checksumming files...")
	matchers := []fs.Matcher{img.GIFMatch, img.JPGMatch, img.PNGMatch}
	for _, d := range cfg.Dirs {
		log.Infof(" - %s", d)
		p, err := fs.NewPath(d, matchers)
		if err != nil {
			log.Error(err)
			continue
		}

		imgPaths, err := p.Find()
		if err != nil {
			log.Error(err)
			continue
		}

		for _, imgPath := range imgPaths {
			log.Info("found image:", imgPath)
			func() {
				fp, err := os.Open(imgPath)
				if err != nil {
					log.Fatal(err)
				}
				defer fp.Close()
				scanStats.ImagesFound++
				hash := md5.New()
				if _, err := io.Copy(hash, fp); err != nil {
					log.Fatal(err)
				}
				md5Sum := fmt.Sprintf("%x", hash.Sum(nil))
				scanStats.FingerPrintCount++
				cfg.Datastore.Put(imgPath, md5Sum)
			}()
		}
	}

	scanStats.Complete()
	log.Info(scanStats)

	return nil
}
