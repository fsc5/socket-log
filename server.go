package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type server struct {
	router        *mux.Router
	socketHandler *wsHandler
	strokeStamps  []KeyStamp
	strokeInput   <-chan KeyStamp
}

type KeyStamp struct {
	Time    time.Time
	IP      string
	Strokes []string
}

func newServer() *server {
	strokeInput := make(chan KeyStamp)
	return &server{
		router:        mux.NewRouter(),
		strokeStamps:  make([]KeyStamp, 0),
		strokeInput:   strokeInput,
		socketHandler: newWsHandler(strokeInput),
	}
}

func (s *server) setUpRoutes() {
	s.router.Handle("/ws", s.socketHandler)
	s.router.HandleFunc("/get", s.outPutStrokes)
	indexTemplate := template.Must(template.ParseFiles("template/index.html"))
	s.router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		indexTemplate.Execute(w, s.strokeStamps)
	})
}

func (s *server) outPutStrokes(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(s.strokeStamps)
}

func (s *server) run() error {
	go func() {
		for stamp := range s.strokeInput {
			s.strokeStamps = append(s.strokeStamps, stamp)
		}
	}()
	return http.ListenAndServe(":80", s.router)
}
