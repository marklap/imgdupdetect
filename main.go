package main

import (
	"flag"

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

}
