package web

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"
	"v2sdl/config"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
)

type WebServer struct {
	config.Service
}

func (w *WebServer) Start(cfg *config.Config) (err error) {
	go w.run(cfg)
	return
}

func (w *WebServer) Name() string                 { return "Web Server" }
func (w *WebServer) Stop()                        {}
func (w *WebServer) Refresh(*config.Config) error { return nil }

func (w *WebServer) run(cfg *config.Config) {
	r := mux.NewRouter()
	r.HandleFunc("/", baseHandler)
	r.HandleFunc("/api/config", func(w http.ResponseWriter, r *http.Request) { writeConfig(buildEnv(cfg), w) })
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

func baseHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Location", "/static/index.html")
	w.WriteHeader(http.StatusFound)
}

func writeConfig(cfg interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	j := json.NewEncoder(w)
	j.SetIndent("", " ")
	j.Encode(cfg)
}

type Environment struct {
	Config     *config.Config
	Interfaces []Interface
	Media      config.Content
}

type Interface struct {
	Name string
	Info string
}

func buildEnv(cfg *config.Config) *Environment {
	displaylist := []Interface{}
	list, _ := net.Interfaces()
	for _, dl := range list {
		addrs, err := dl.Addrs()
		if err != nil || len(addrs) == 0 {
			continue
		}

		displays := []string{}
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					displays = append(displays, ipnet.IP.String())
				}
			}
		}
		if len(displays) == 0 {
			continue
		}

		info := fmt.Sprintf("%s %s", dl.Name, strings.Join(displays, " "))
		displaylist = append(displaylist, Interface{
			Name: dl.Name,
			Info: info,
		})
	}

	return &Environment{
		Config:     cfg,
		Interfaces: displaylist,
		Media:      config.Media,
	}
}
