package main

import (
	"flag"
	"os"

	log "github.com/sirupsen/logrus"
)

func main() {
	var debug = flag.Bool("debug", false, "turn debug logging on")
	var listen = flag.String("listen", "127.0.0.1:8228", "interface to listen for web user interface")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

	dirs := flag.Args()
	err := httpServer(*listen, dirs)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	os.Exit(0)
}
