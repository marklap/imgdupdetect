package cli

import (
	"github.com/marklap/imgdupdetect/datastore"
	"github.com/marklap/imgdupdetect/fs"
	"github.com/marklap/imgdupdetect/img"
	"github.com/marklap/imgdupdetect/stats"

	log "github.com/sirupsen/logrus"
)

// Config is the CLI config
type Config struct {

	// Dirs is the directories to scan for duplicates
	Dirs []string
	// Datastore is the datastore
	Datastore *datastore.Datastore
	// FingerPrintCol is the name of the collection to use for fingerprints
	FingerPrintCol string
}

// Run runs the specified command
func Run(cfg Config, cmd string) error {
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
