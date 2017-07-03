package cli

import (
	"sync"

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

// msg is a queue message
type msg struct {
	img *img.Image
	fp  []byte
}

// imageFingerPrint FingerPrints the Images
func imageFingerPrint(images <-chan *img.Image, fpQueue chan<- *msg, scanStats *stats.ScanStats) {
	for {
		select {
		case image, ok := <-images:
			if !ok {
				// channel is closed
				return
			}
			fp, err := image.FingerPrint()
			if err != nil {
				log.Error(err)
				return
			}
			scanStats.FingerPrintCountIncr()

			fpQueue <- &msg{image, fp}
		}
	}
}

// recordFingerPrint

// Run runs the specified command
func Run(cfg *Config, cmd string) error {
	log.Info("looking for duplicates...")

	scanStats := stats.NewScanStats()

	images := make(chan *img.Image)
	fpQueue := make(chan *msg)

	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU(); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			imageFingerPrint(images, fpQueue, scanStats)
		}()
	}

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

			images <- i

			// scanStats.ImagesFound++

			// fp, err := i.FingerPrint()
			// if err != nil {
			// 	log.Error(err)
			// 	continue
			// }
			// scanStats.FingerPrintCount++

			// meta := map[string][]byte{
			// 	"size":   i.SizeByteSlice(),
			// 	"height": i.HeightByteSlice(),
			// 	"width":  i.WidthByteSlice(),
			// }

			// err = cfg.Datastore.Add(cfg.FingerPrintCol, fp, imgPath, meta)
			// if err != nil {
			// 	log.Error(err)
			// 	continue
			// }
		}
	}

	go func() {
		for m := range fpQueue {
			meta := map[string][]byte{
				"size":   m.img.SizeByteSlice(),
				"height": m.img.HeightByteSlice(),
				"width":  m.img.WidthByteSlice(),
			}

			err := cfg.Datastore.Add(cfg.FingerPrintCol, m.fp, m.img.Path, meta)
			if err != nil {
				log.Error(err)
				continue
			}
		}
	}()

	wg.Wait()
	close(fpQueue)

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
