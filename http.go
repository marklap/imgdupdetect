package main

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"golang.org/x/net/websocket"
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
<head>
<title>Image Duplicate Detector</title>
</head>
<body>
%s
<a href="javascript:WebSocketTest()">WebSocket Test</a>
<script>
	function WebSocketTest() {
		if ("WebSocket" in window) {
			alert("WebSocket is supported by your Browser!");

			// Let us open a web socket
			var ws = new WebSocket("ws://localhost:8228/ws");

			ws.onopen = function() {
				// Web Socket is connected, send data using send()
				ws.send("Message to send");
				alert("Message is sent...");
			};

			ws.onmessage = function (evt) {
				var received_msg = evt.data;
				alert("Message is received...");
			};

			ws.onclose = function() {
				// websocket is closed.
				alert("Connection is closed...");
			};
		} else {
			// The browser doesn't support WebSocket
			alert("WebSocket NOT supported by your Browser!");
		}
	}
</script>
</body>
</html>`

func wsHandler(ws *websocket.Conn) {

	after10Sec := time.Now().Add(10 * time.Second)
	everySec := time.Tick(time.Second)
	for now := range everySec {
		ws.Write([]byte(fmt.Sprintf("%s", now)))
		if now.After(after10Sec) {
			break
		}
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	cfg := r.Context().Value(ctxKeyConfig).(httpConfig)
	w.Header().Set("Access-Control-Allow-Origin", "*")
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
	log.Infof("starting http server on %s", cfg.listen)
	http.Handle("/ws", websocket.Handler(wsHandler))
	http.Handle("/", http.Handler(withConfig(cfg, indexHandler)))
	return http.ListenAndServe(cfg.listen, nil)
}
