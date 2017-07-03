package cli

import (
	"github.com/marklap/imgdupdetect/datastore"
	"github.com/marklap/imgdupdetect/fs"
	"github.com/marklap/imgdupdetect/img"
	"github.com/marklap/imgdupdetect/stats"

	log "github.com/sirupsen/logrus"
	"runtime"
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

// processImage FingerPrints the Images
func processImage(cfg Config, path string, scanStats *stats.ScanStats) {
	i, err := img.NewImage(path)
	if err != nil {
		log.Error(err)
		return
	}
	scanStats.ImagesFoundIncr()

	fp, err := i.FingerPrint()
	if err != nil {
		log.Error(err)
		return
	}
	scanStats.FingerPrintCountIncr()

	meta := map[string][]byte{
		"size":   i.SizeByteSlice(),
		"height": i.HeightByteSlice(),
		"width":  i.WidthByteSlice(),
	}

	err = cfg.Datastore.Add(cfg.FingerPrintCol, fp, i.Path, meta)
	if err != nil {
		log.Error(err)
		return
	}
}

// worker is kicks off the imageProcess for each image path
func worker(cfg Config, paths <-chan string, scanStats *stats.ScanStats) {
	for p := range paths {
		processImage(cfg, p, scanStats)
	}
}

// Run runs the specified command
func Run(cfg Config, cmd string) error {
	log.Info("looking for duplicates...")

	scanStats := stats.NewScanStats()

	imagePaths := make(chan string)

	for i := 0; i < runtime.NumCPU(); i++ {
		go worker(cfg, imagePaths, scanStats)
	}

	matchers := []fs.Matcher{img.GIFMatch, img.JPGMatch, img.PNGMatch}
	for _, d := range cfg.Dirs {
		log.Infof(" - %s", d)
		p, err := fs.NewPath(d, matchers)
		if err != nil {
			log.Error(err)
			continue
		}

		paths, err := p.Find()
		if err != nil {
			log.Error(err)
			continue
		}

		for _, path := range paths {
			imagePaths <- path
		}
	}

	close(imagePaths)

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
