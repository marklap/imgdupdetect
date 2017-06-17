package main

import (
	"flag"
	// "image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	// "io/ioutil"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

func main() {
	var (
		debug   = flag.Bool("debug", false, "turn debug logging on")
		root    = flag.String("root", ".", "root directory")
		pattern = flag.String("pattern", "*.jpeg", "pattern to match for files")
	)
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	err := filepath.Walk(*root, func(path string, info os.FileInfo, err error) error {
		if matched, err := filepath.Match(*pattern, filepath.Base(path)); err == nil {
			if matched {
				log.Info("matches:", path)
			} else {
				log.Info("not this one:", path)
			}
		} else {
			log.Error(err)
		}
		return nil
	})
	if err != nil {
		log.Error(err)
	}

}
