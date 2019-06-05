package main

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/marklap/imgdupdetect/cli"
	"github.com/marklap/imgdupdetect/datastore"
	"github.com/marklap/imgdupdetect/ui"

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
	var datastorePath = flag.String("datastore", filepath.Join(here, "imgdd.ds"), "path where the datastore should be saved")
	var debug = flag.Bool("debug", false, "turn debug logging on")
	var listen = flag.String("listen", "127.0.0.1:8228", "interface to listen for web user interface")
	var serveHTTP = flag.Bool("ui", false, "start an http server at `listen`")
	var relocateFrom = flag.String("relo-from", "", "relocate images from path")
	var relocateTo = flag.String("relo-to", "", "relocate images to path")
	var dumpDB = flag.Bool("dump", false, "dump contents of database and exit")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	if (len(*relocateFrom) > 0 && *relocateTo == "") || (len(*relocateTo) > 0 && *relocateFrom == "") {
		log.Error("must specify relocate from and relocate to")
		os.Exit(1)
	}

	var dirs []string
	if len(*relocateFrom) > 0 && len(*relocateTo) > 0 {

		if _, err := os.Stat(*relocateFrom); os.IsNotExist(err) {

			log.Error("relo-from directory does not exist: ", *relocateFrom)
			os.Exit(1)
		}
		if _, err := os.Stat(*relocateTo); os.IsNotExist(err) {
			log.Error("relo-to directory does not exist: ", *relocateTo)
			os.Exit(1)
		}
		if *relocateFrom == *relocateTo {
			log.Error("relo-from and relo-to must not be the same directory")
			os.Exit(1)
		}
	} else if dirs = flag.Args(); len(dirs) == 0 {
		log.Error("no directories specified")
		os.Exit(1)
	}

	log.Debugf("static dir: %s", *static)

	ds, err := datastore.Open(datastore.Config{Path: *datastorePath})
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	defer ds.Close()

	if *dumpDB {
		ds.Dump(os.Stdout)
		os.Exit(0)
	}

	if *serveHTTP {
		err = ui.Serve(ui.Config{
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
	} else if len(*relocateFrom) > 0 && len(*relocateTo) > 0 {
		err = cli.ReloRun(cli.ReloConfig{
			From: *relocateFrom,
			To:   *relocateTo,
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
		err = cli.MD5Sum(cli.DupeDetectConfig{
			Dirs:           dirs,
			Datastore:      ds,
			FingerPrintCol: "md5sum",
		}, cmd)
		if err != nil {
			log.Error(err)
			os.Exit(1)
		}
	}
	os.Exit(0)
}
