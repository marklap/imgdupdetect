package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
)

// httpConfig is the http server Config
type httpConfig struct {
	// dirs are the directories to look in for duplicates
	dirs []string
	// listen is the interface an port to listen on
	listen string
}

// ctxKey is the context type key
type ctxKey string

// ctxKeyConfig specifies the context key for the config
const ctxKeyConfig = ctxKey("ctxConfig")

var indexTmpl = `<html>
<head><title>Image Duplicate Detector</title></head>
<body>%s</body>
</html>`

func indexHandler(w http.ResponseWriter, r *http.Request) {
	cfg := r.Context().Value(ctxKeyConfig).(httpConfig)
	w.Write([]byte(fmt.Sprintf(indexTmpl, strings.Join(cfg.dirs, "<br />"))))
}

// withConfig wraps an http.HandlerFunc with a context that contains all the http server config
func withConfig(cfg httpConfig, f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxKeyConfig, cfg)
		f(w, r.WithContext(ctx))
	})
}

func httpServer(cfg httpConfig) error {
	http.Handle("/", http.Handler(withConfig(cfg, indexHandler)))
	return http.ListenAndServe(cfg.listen, nil)
}
