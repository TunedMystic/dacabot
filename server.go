package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func newServer() *Server {
	s := Server{}
	fmt.Println("[setup] templates")
	s.templates = template.Must(template.ParseGlob("templates/*.html"))

	fmt.Println("[setup] database")
	s.db = newDB()
	s.db.createTables()

	fmt.Println("[setup] router")
	s.router = mux.NewRouter()
	return &s
}

type Server struct {
	templates *template.Template
	router    *mux.Router
	db        *ServerDB
}

func (s *Server) indexHandler(w http.ResponseWriter, r *http.Request) {
	s.templates.ExecuteTemplate(w, "index", nil)
}

func (s *Server) articlesHandler(w http.ResponseWriter, r *http.Request) {
	articles := s.db.fetchArticles()
	data := TemplateContext{Articles: articles}
	s.templates.ExecuteTemplate(w, "articles", data)
}

func (s *Server) healthView(w http.ResponseWriter, r *http.Request) {
	err := s.db.checkhealth()
	if err != nil {
		http.Error(w, "Database health check failed", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, "ok")
}

func (s *Server) run() {
	s.router.HandleFunc("/", s.indexHandler).Methods("GET")
	s.router.HandleFunc("/articles", s.articlesHandler).Methods("GET")
	s.router.HandleFunc("/health", s.healthView).Methods("GET")
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	fmt.Println("[run] starting Server on port 8000...")
	log.Fatal(http.ListenAndServe("0.0.0.0:8000", s.router))
}
