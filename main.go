package main

import (
	"flag"

	log "github.com/sirupsen/logrus"
)

func main() {
	var debug = flag.Bool("debug", false, "turn debug logging on")
	flag.Parse()

	if *debug {
		log.SetLevel(log.DebugLevel)
	} else {
		log.SetLevel(log.InfoLevel)
	}

}
