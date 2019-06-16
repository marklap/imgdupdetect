package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/marklap/imgdupdetect/cli"
	"github.com/marklap/imgdupdetect/datastore"

	log "github.com/sirupsen/logrus"
)

const (
	statCollection = "stat"
)

func main() {
	here, err := filepath.Abs(".")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	var datastorePath = flag.String("datastore", filepath.Join(here, "imgdd.sqlite"), "path where the datastore should be saved")
	var debug = flag.Bool("debug", false, "turn debug logging on")
	// var dumpDB = flag.Bool("dump", false, "dump contents of database and exit")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	var dirs []string
	if dirs = flag.Args(); len(dirs) == 0 {
		log.Error("no directories specified")
		os.Exit(1)
	}

	ds, err := datastore.Open(datastore.Config{Path: *datastorePath})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer ds.Close()

	err = cli.MD5Sum(cli.DupeDetectConfig{
		Dirs:      dirs,
		Datastore: ds,
	})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	os.Exit(0)
}
