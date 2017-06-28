package main

import (
	"net/http"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

func httpServer(listen string, dirs []string) error {
	http.HandleFunc("/", indexHandler)
	return http.ListenAndServe(listen, nil)
}
