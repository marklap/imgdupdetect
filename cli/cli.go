package cli

import (
	"crypto/md5"
	"crypto/rand"
	"fmt"
	"io"
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

		newPath, isDup, sz := uniqueDestPath(
			filepath.Join(cfg.To, dt.Format("2006-01-02")),
			fd,
		)

		fmt.Printf("\"%s\",\"%s\",\"%t\",%d\n", path, newPath, isDup, sz)

		err = ensureDir(newPath)
		if err != nil {
			log.Error(err)
			return err
		}

		destFile, err := os.Create(newPath)
		if err != nil {
			log.Error(err)
			return err
		}

		_, err = io.Copy(fd, destFile)
		if err != nil {
			log.Error(err)
			return err
		}
	}

	return nil
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

func uniqueDestPath(dstPath string, src *os.File) (string, bool, int64) {
	srcBase := filepath.Base(src.Name())
	srcStat, _ := src.Stat()
	srcSize := srcStat.Size()
	newPath := filepath.Join(dstPath, srcBase)

	if _, err := os.Stat(newPath); os.IsNotExist(err) {
		return newPath, false, srcSize
	}

	fileParts := strings.SplitN(srcBase, ".", 2)
	fileName, fileExt := fileParts[0], fileParts[1]

	for i := 1; i < 100; i++ {
		newPath := filepath.Join(dstPath, fmt.Sprintf("%s.%03d.%s", fileName, i, fileExt))
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return newPath, true, srcSize
		}
	}

	return "", false, 0
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

func MD5Sum(cfg DupeDetectConfig, cmd string) error {
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
			log.Info("found image:", imgPath)
			func() {
				fp, err := os.Open(imgPath)
				if err != nil {
					log.Fatal(err)
				}
				defer fp.Close()
				hash := md5.New()
				if _, err := io.Copy(hash, fp); err != nil {
					log.Fatal(err)
				}
				log.Info(fmt.Sprintf("\tmd5sum: %x", hash.Sum(nil)))
			}()
		}
	}

	scanStats.Complete()
	log.Info(scanStats)

	return nil
}
