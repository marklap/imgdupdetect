package cli

import (
	"crypto/rand"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/marklap/imgdupdetect/datastore"
	"github.com/marklap/imgdupdetect/fs"
	"github.com/marklap/imgdupdetect/img"
	"github.com/marklap/imgdupdetect/stats"

	"github.com/rwcarlsen/goexif/exif"
	"github.com/rwcarlsen/goexif/mknote"

	log "github.com/sirupsen/logrus"
)

func uuid() string {
	buf := make([]byte, 6)
	_, err := rand.Read(buf)
	if err != nil {
		return ""
	}
	id := fmt.Sprintf("%012x", buf)
	return id
}

// DupeDetectConfig is the duplicate detector CLI config
type DupeDetectConfig struct {

	// Dirs is the directories to scan for duplicates
	Dirs []string
	// Datastore is the datastore
	Datastore *datastore.Datastore
	// FingerPrintCol is the name of the collection to use for fingerprints
	FingerPrintCol string
}

// ReloConfig is the relocation CLI config
type ReloConfig struct {

	// From is the directory of the source images
	From string
	// To is the target directory
	To string
}

// ReloRun runs the relocation function
func ReloRun(cfg ReloConfig) error {
	log.Debug("relo from: ", cfg.From)
	log.Debug("relo to: ", cfg.To)

	p, err := fs.NewPath(cfg.From, []fs.Matcher{img.TIFFMatch, img.JPGMatch})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	imgPaths, err := p.Find()
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	exif.RegisterParsers(mknote.All...)

	for _, path := range imgPaths {
		var err error
		// var i *img.Image
		// i, err = img.NewImage(path)
		// if err != nil {
		// 	log.Error(err)
		// }

		fd, err := os.Open(path)
		defer fd.Close()

		if err != nil {
			log.Error(err)
			return err
		}

		ximg, err := exif.Decode(fd)
		if err != nil {
			log.Error(err)
			return err
		}

		dt, _ := ximg.DateTime()
		if dt.IsZero() {
			log.Error("not a date")
			continue
		}

		base := filepath.Base(path)
		parts := strings.SplitN(base, ".", 2)

		fmt.Println(
			path,
			filepath.Join(
				cfg.To,
				dt.Format("2006-01-02"),
				strings.ToLower(strings.Join([]string{parts[0], uuid(), parts[1]}, ".")),
			),
		)
	}

	return nil
}

// DupeDetectRun runs the duplicate detect function
func DupeDetectRun(cfg DupeDetectConfig, cmd string) error {
	scanStats := stats.NewScanStats()

	log.Info("looking for duplicates...")
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
			i, err := img.NewImage(imgPath)
			if err != nil {
				log.Error(err)
				continue
			}
			scanStats.ImagesFound++

			meta := map[string][]byte{
				"size":   i.SizeByteSlice(),
				"height": i.HeightByteSlice(),
				"width":  i.WidthByteSlice(),
			}

			fp, err := i.FingerPrint()
			if err != nil {
				log.Error(err)
				continue
			}
			scanStats.FingerPrintCount++

			err = cfg.Datastore.Add(cfg.FingerPrintCol, fp, imgPath, meta)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}

	fps := cfg.Datastore.GetFingerPrints(cfg.FingerPrintCol)
	for _, fp := range fps {
		ims := cfg.Datastore.GetImages(cfg.FingerPrintCol, fp)
		imsLen := len(ims)
		if imsLen > 1 {
			scanStats.DuplicatesFound += imsLen - 1 // we don't count the original
			log.Info("found duplicates:")
			for _, i := range ims {
				log.Infof("  - %s", i)
			}
		}
	}

	scanStats.Complete()
	log.Info(scanStats)

	return nil
}
