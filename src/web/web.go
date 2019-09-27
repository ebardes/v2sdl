package web

import (
	"net/http"
	"time"
	"v2sdl/config"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type WebServer struct {
}

func (w *WebServer) Start(cfg *config.Config) {
	go w.run(cfg)
}

func (w *WebServer) run(cfg *config.Config) {
	r := mux.NewRouter()
	r.HandleFunc("/", BaseHandler)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	r.PathPrefix("/content/").Handler(http.StripPrefix("/content/", http.FileServer(http.Dir(cfg.Storage))))

	srv := &http.Server{
		Handler:      r,
		Addr:         ":8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	srv.SetKeepAlivesEnabled(true)

	log.Info().Msgf("Launching web server at: %s", srv.Addr)
	log.Error().Msgf("Error launching web server: %v", srv.ListenAndServe())
}

func BaseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "/static/index.html")
	w.WriteHeader(http.StatusFound)
}
