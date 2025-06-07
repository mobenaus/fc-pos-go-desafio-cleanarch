package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router         chi.Router
	Handlers       map[string]http.HandlerFunc
	POSTHandlers   map[string]http.HandlerFunc
	GETHandlers    map[string]http.HandlerFunc
	PUTHandlers    map[string]http.HandlerFunc
	DELETEHandlers map[string]http.HandlerFunc
	WebServerPort  string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:         chi.NewRouter(),
		Handlers:       make(map[string]http.HandlerFunc),
		POSTHandlers:   make(map[string]http.HandlerFunc),
		GETHandlers:    make(map[string]http.HandlerFunc),
		PUTHandlers:    make(map[string]http.HandlerFunc),
		DELETEHandlers: make(map[string]http.HandlerFunc),
		WebServerPort:  serverPort,
	}
}

func (s *WebServer) AddHandler(path string, handler http.HandlerFunc) {
	s.Handlers[path] = handler
}
func (s *WebServer) AddPOSTHandler(path string, handler http.HandlerFunc) {
	s.POSTHandlers[path] = handler
}
func (s *WebServer) AddGETHandler(path string, handler http.HandlerFunc) {
	s.GETHandlers[path] = handler
}
func (s *WebServer) AddPUTHandler(path string, handler http.HandlerFunc) {
	s.PUTHandlers[path] = handler
}
func (s *WebServer) AddDELETEHandler(path string, handler http.HandlerFunc) {
	s.DELETEHandlers[path] = handler
}

// loop through the handlers and add them to the router
// register middeleware logger
// start the server
func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	for path, handler := range s.Handlers {
		s.Router.Handle(path, handler)
	}
	for path, handler := range s.POSTHandlers {
		s.Router.Post(path, handler)
	}
	for path, handler := range s.GETHandlers {
		s.Router.Get(path, handler)
	}
	for path, handler := range s.PUTHandlers {
		s.Router.Put(path, handler)
	}
	for path, handler := range s.DELETEHandlers {
		s.Router.Delete(path, handler)
	}
	http.ListenAndServe(s.WebServerPort, s.Router)
}
