package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/marklap/imgdupdetect/cli"
	"github.com/marklap/imgdupdetect/datastore"
	"github.com/marklap/imgdupdetect/gui"

	log "github.com/sirupsen/logrus"
)

const (
	statCollection        = "stat"
	fingerPrintCollection = "fingerprint"
)

func main() {
	here, err := filepath.Abs(".")
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}

	var static = flag.String("static", filepath.Join(here, "static"), "path where static files (html, css, js) are located")
	var datastorePath = flag.String("datastore", filepath.Join(here, "imgdupdetect.ds"), "path where the datastore should be saved")
	var debug = flag.Bool("debug", false, "turn debug logging on")
	var listen = flag.String("listen", "127.0.0.1:8228", "interface to listen for web user interface")
	var serveHTTP = flag.Bool("gui", false, "start an http server at `listen` for a gui")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	log.Debugf("static dir: %s", *static)

	dirs := flag.Args()
	if len(dirs) == 0 {
		log.Error("no directories specified")
		os.Exit(1)
	}

	ds, err := datastore.Open(datastore.Config{Path: *datastorePath})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer ds.Close()

	if *serveHTTP {
		err = gui.Serve(gui.Config{
			Dirs:           dirs,
			Listen:         *listen,
			Static:         *static,
			Datastore:      ds,
			FingerPrintCol: fingerPrintCollection,
		})
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	} else {
		var cmd = "fingerprint"
		if len(dirs) > 1 {
			switch dirs[0] {
			case "clear":
				cmd = "clear"
				dirs = dirs[1:]
			case "fingerprint":
				dirs = dirs[1:]
			default:
			}
		}
		err = cli.Run(cli.Config{
			Dirs:           dirs,
			Datastore:      ds,
			FingerPrintCol: fingerPrintCollection,
		}, cmd)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
