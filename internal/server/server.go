package server

import (
	"io/fs"
	"net/http"
)

type Server struct {
	webFS     fs.FS
	exampleFS fs.FS
	version   string
}

func New(webFS, exampleFS fs.FS, version string) *Server {
	return &Server{webFS: webFS, exampleFS: exampleFS, version: version}
}

func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.handleHealth)
	mux.HandleFunc("/api/health", s.handleHealth)
	mux.HandleFunc("/api/conc", s.handleConc)
	mux.HandleFunc("/api/axis", s.handleAxis)
	mux.HandleFunc("/api/version", s.handleVersion)
	mux.Handle("/example/", http.StripPrefix("/example/", http.FileServer(http.FS(s.exampleFS))))
	mux.Handle("/", http.FileServer(http.FS(s.webFS)))
	return mux
}

func (s *Server) Version() string { return s.version }
