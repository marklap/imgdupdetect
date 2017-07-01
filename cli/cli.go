package cli

import (
	"github.com/marklap/imgdupdetect/datastore"
	"github.com/marklap/imgdupdetect/fs"
	"github.com/marklap/imgdupdetect/img"

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
	log.Info("looking for duplicates...")
	matchers := []fs.Matcher{img.GIFMatch, img.JPGMatch, img.PNGMatch}
	for _, d := range cfg.Dirs {
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
		if len(ims) > 1 {
			log.Info("found duplicates:")
			for _, i := range ims {
				log.Infof("  - %s", i)
			}
		}
	}

	return nil
}
