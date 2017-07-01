package gui

import (
	"context"
	"net/http"
	"os"
	"path/filepath"

	"github.com/marklap/imgdupdetect/datastore"
	// "github.com/marklap/imgdupdetect/fs"
	// "github.com/marklap/imgdupdetect/img"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
)

const (
	staticDir = "static"
	indexHTML = "index.html"
)

// Config is the http server Config
type Config struct {
	// dirs are the directories to look in for duplicates
	Dirs []string
	// listen is the interface an port to listen on
	Listen string
	// static is the directeory where static (css, js, html) files are located
	Static string
	// Datastore is the datastore
	Datastore *datastore.Datastore
	// FingerPrintCol is the name of the collection to use for fingerprints
	FingerPrintCol string
}

// ctxKey is the context type key
type ctxKey string

// ctxKeyConfig specifies the context key for the config
const ctxKeyConfig = ctxKey("ctxConfig")

func wsHandler(ws *websocket.Conn) {
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	cfg := r.Context().Value(ctxKeyConfig).(Config)
	indexFilePath := filepath.Join(cfg.Static, "html", indexHTML)
	index, err := os.Open(indexFilePath)
	if err != nil {
		log.Errorf("failed to open index html flie: %s", indexFilePath)
		w.WriteHeader(http.StatusInternalServerError)
	}
	indexStat, err := index.Stat()
	if err != nil {
		log.Errorf("failed to stat index html file: %s", indexFilePath)
		w.WriteHeader(http.StatusInternalServerError)
	}
	defer index.Close()

	w.Header().Set("Access-Control-Allow-Origin", "*")
	http.ServeContent(w, r, indexFilePath, indexStat.ModTime(), index)
}

// withConfig wraps an http.HandlerFunc with a context that contains all the http server config
func withConfig(cfg Config, f http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := context.WithValue(r.Context(), ctxKeyConfig, cfg)
		f(w, r.WithContext(ctx))
	})
}

// Serve initiates the HTTP server
func Serve(cfg Config) error {
	log.Infof("starting http server on %s", cfg.Listen)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir(cfg.Static))))
	http.Handle("/ws", websocket.Handler(wsHandler))
	http.Handle("/", http.Handler(withConfig(cfg, indexHandler)))
	return http.ListenAndServe(cfg.Listen, nil)
}
